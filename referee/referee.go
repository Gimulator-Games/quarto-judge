package referee

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/Gimulator-Games/quarto-judge/types"
	"github.com/Gimulator/client-go"
	"github.com/sirupsen/logrus"
)

type Referee struct {
	cli

	log    *logrus.Entry
	board  types.Board
	ch     chan client.Object
	p1, p2 types.Player
	roomID string
}

func NewReferee(roomID string) (Referee, error) {
	ch := make(chan client.Object)

	controller, err := newCli(ch, roomID)
	if err != nil {
		return Referee{}, err
	}

	return Referee{
		cli:    controller,
		ch:     ch,
		log:    logrus.WithField("entity", "referee"),
		roomID: roomID,
	}, nil
}

func (r *Referee) Start() error {
	r.log.Info("starting referee")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	r.log.Info("starting to receipt playes")
	o1, o2, n := r.receiptPlayers(ctx)
	if n != 2 {
		return fmt.Errorf("could not find two anget for starting the game")
	}
	r.handlePlayers(o1, o2)
	r.log.WithField("player1", r.p1).WithField("player2", r.p2).Debug("players are receipted")

	r.board = types.NewBoard(r.p1.Name)

	r.log.Info("starting to set initial board")
	if err := r.setBoard(r.board); err != nil {
		r.log.WithError(err).Error("could not set the initial board of the game")
		return err
	}

	if err := r.listen(); err != nil {
		r.log.WithError(err).Error("could not listen")
		return err
	}
	r.log.Info("end of game")

	return nil
}

func (r *Referee) listen() error {
	for {
		r.log.Info("starting to listen")
		timer := time.NewTimer(time.Second * 3)
		select {
		case <-timer.C:
			r.log.WithField("turn", r.board.Turn).Debug("turn's time is over")
			return r.endOfGameWithTimeout()
		case obj := <-r.ch:
			if obj.Key.Type != actionType {
				r.log.Debug("invalid type of object, type should be action")
				continue
			}

			if err := r.judge(obj); err != nil {
				r.log.WithError(err).Error("could not judge the recieved object")
			}
		}
	}
}

func (r *Referee) judge(obj client.Object) error {
	r.log.WithField("object", obj).Info("starting to judge recieved object")

	if err, isLosing := r.checkTurn(obj); isLosing {
		return r.endOfGameWithLoser(obj.Meta.Owner)
	} else if err != nil {
		return err
	}

	action := types.Action{}
	if err := json.Unmarshal([]byte(obj.Value), &action); err != nil {
		return err
	}
	r.log.WithField("action", action).Info("starting to judge recieved action")

	if !r.validateAction(action) {
		return r.endOfGameWithLoser(obj.Meta.Owner)
	}
	r.update(action)

	if r.isWinState() {
		r.endOfGameWithWinner(obj.Meta.Owner)
		return nil
	}

	if r.isTieState() {
		r.endOfGameTie()
		return nil
	}

	r.changeTurn()

	if err := r.setBoard(r.board); err != nil {
		return err
	}

	return nil
}

func (r *Referee) checkTurn(obj client.Object) (error, bool) {
	if obj.Key.Name != r.board.Turn {
		return fmt.Errorf("invalid turn move: it's %s's turn", r.board.Turn), false
	}
	if obj.Key.Name == r.p1.Name && obj.Meta.Owner != r.p1.Id {
		return fmt.Errorf("invalid owner: incoming object has invalid owner = '%s' in meta", obj.Meta.Owner), true
	}
	if obj.Key.Name == r.p2.Name && obj.Meta.Owner != r.p2.Id {
		return fmt.Errorf("invalid owner: incoming object has invalid owner = '%s' in meta", obj.Meta.Owner), true
	}
	return nil, false
}

func (r *Referee) update(action types.Action) {
	for i := range r.board.Positions {
		pos := &r.board.Positions[i]
		if action.X == pos.X && action.Y == pos.Y {
			pos.PieceID = r.board.Picked
		}
	}

	r.board.Picked = action.Picked
}

func (r *Referee) changeTurn() {
	switch r.board.Turn {
	case r.p1.Name:
		r.board.Turn = r.p2.Name
	case r.p2.Name:
		r.board.Turn = r.p1.Name
	default:
		r.log.WithFields(logrus.Fields{
			"p1":   r.p1.Name,
			"p2":   r.p2.Name,
			"turn": r.board.Turn,
		}).Fatal("could not change the board's turn")
		r.board.Turn = ""
	}
}

func (r *Referee) validateAction(action types.Action) bool {
	r.log.Info("starting to validate action")

	if action.X < 1 || action.X > 4 || action.Y < 1 || action.Y > 4 {
		fmt.Println("========================0")
		return false
	}

	if action.Picked == r.board.Picked {
		fmt.Println("===========6=============1")
		return false
	}

	for _, pos := range r.board.Positions {
		if pos.X == action.X && pos.Y == action.Y && pos.PieceID != 0 {
			fmt.Println("========================2")
			return false
		}

		if action.Picked == pos.PieceID {
			fmt.Println("========================3")
			return false
		}
	}

	if action.Picked < 1 || action.Picked > 16 {
		fmt.Println("========================4")
		return false
	}

	r.log.Debug("action is valid")
	return true
}

func (r *Referee) handlePlayers(o1, o2 *client.Object) {
	rnd := rand.Intn(2)
	if rnd == 0 {
		r.p1 = types.Player{
			Name: o1.Key.Name,
			Id:   o1.Meta.Owner,
		}
		r.p2 = types.Player{
			Name: o2.Key.Name,
			Id:   o2.Meta.Owner,
		}
	} else {
		r.p1 = types.Player{
			Name: o2.Key.Name,
			Id:   o2.Meta.Owner,
		}
		r.p2 = types.Player{
			Name: o1.Key.Name,
			Id:   o1.Meta.Owner,
		}
	}
}

func (r *Referee) endOfGameWithTimeout() error {
	var err error
	if r.board.Turn == r.p1.Name {
		err = r.endOfGameWithLoser(r.p1.Id)
	} else {
		err = r.endOfGameWithLoser(r.p2.Id)
	}
	return err
}

func (r *Referee) endOfGameWithWinner(winner string) error {
	r.log.WithField("winner", winner).Info("end of game")
	if r.p1.Id == winner {
		return r.endOfGame(r.p1.Id, r.p2.Id, 3, 0)
	}
	return r.endOfGame(r.p2.Id, r.p1.Id, 3, 0)
}

func (r *Referee) endOfGameWithLoser(loser string) error {
	r.log.WithField("loser", loser).Info("end of game")
	if r.p1.Id == loser {
		return r.endOfGame(r.p2.Id, r.p1.Id, 3, 0)
	}
	return r.endOfGame(r.p1.Id, r.p2.Id, 3, 0)
}

func (r *Referee) endOfGameTie() error {
	r.log.Info("end of game in tie state")
	return r.endOfGame(r.p1.Id, r.p2.Id, 1, 1)
}

func (r *Referee) endOfGame(p1, p2 string, s1, s2 int) error {
	r.board.Turn = ""
	if err := r.setBoard(r.board); err != nil {
		return err
	}

	scores := make(map[string]map[int]int)
	scores[p1] = map[int]int{1: s1}
	scores[p2] = map[int]int{1: s2}

	res := result{
		RoomID:  r.roomID,
		Status:  "SUCCESS",
		Message: "",
		Scores:  scores,
	}

	return r.setEndOfGame(res)
}

func (r *Referee) isWinState() bool {
	b := [5][5]int{}
	for _, pos := range r.board.Positions {
		if pos.PieceID == 0 {
			b[pos.X][pos.Y] = 0
		} else {
			b[pos.X][pos.Y] = r.board.Pieces[pos.PieceID].Code()
		}
	}

	andDiamM := 255
	andDiamS := 255

	for i := 1; i <= 4; i++ {
		andRow := 255
		andCol := 255
		for j := 1; j <= 4; j++ {
			andRow &= b[i][j]
			andCol &= b[j][i]
		}

		if andRow != 0 || andCol != 0 {
			r.log.WithField("index", i).Debug("row or col cause to win")
			return true
		}

		andDiamM &= b[i][i]
		andDiamS &= b[i][5-i]
	}

	if andDiamM != 0 || andDiamS != 0 {
		r.log.Debug("Diams cause to win")
		return true
	}
	return false
}

func (r *Referee) isTieState() bool {
	for _, pos := range r.board.Positions {
		if pos.PieceID == 0 {
			return false
		}
	}

	return true
}
