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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	o1, o2, n := r.receiptPlayers(ctx)
	if n != 2 {
		r.log.WithField("object1", o1).WithField("object2", o2).Error("could not find two anget for starting the game")
		return fmt.Errorf("could not find two anget for starting the game")
	}
	r.handlePlayers(o1, o2)

	r.board = types.NewBoard(r.p1.Name)

	if err := r.setBoard(r.board); err != nil {
		r.log.WithError(err).Error("could not set the initial board of the game")
		return err
	}

	r.log.Info("starting to listen")
	if err := r.listen(); err != nil {
		r.log.WithError(err).Error("could not listen")
		return err
	}
	r.log.Info("end of game")

	return nil
}

func (r *Referee) listen() error {
	for {
		timer := time.NewTimer(time.Second * 3)
		select {
		case <-timer.C:
			r.log.Debug("turn's time is over")
			r.endOfGameWithLoser(r.board.Turn)
			return nil
		case obj := <-r.ch:
			if obj.Key.Type != actionType {
				r.log.Debug("invalid type of object, type should be action")
				continue
			}

			if r.judge(obj) {
				break
			}
		}
	}
}

func (r *Referee) judge(obj client.Object) bool {
	if obj.Key.Name != r.board.Turn {
		r.log.Info("invalid turn move: it's %s's turn", r.board.Turn)
		return false
	}
	if obj.Key.Name == r.p1.Name && obj.Meta.Owner != r.p1.Id {
		r.log.Info("invalid owner: incoming object has invalid owner = '%s' in meta", obj.Meta.Owner)
		r.endOfGameWithLoser(r.p1.Name)
		return true
	}
	if obj.Key.Name == r.p2.Name && obj.Meta.Owner != r.p2.Id {
		r.log.Info("invalid owner: incoming object has invalid owner = '%s' in meta", obj.Meta.Owner)
		r.endOfGameWithLoser(r.p1.Name)
		return true
	}

	log := r.log.WithField("turn", obj.Key.Name)

	bytes, ok := obj.Value.([]byte)
	if !ok {
		log.Error("could not cast value to bytes")
		return false
	}

	action := types.Action{}
	if err := json.Unmarshal(bytes, &action); err != nil {
		log.Error("could not unmarshal action")
		return false
	}

	if !r.validateAction(action) {
		log.Info("invalid action")
		r.endOfGameWithLoser(r.board.Turn)
		return true
	}
	r.update(action)

	if r.isWinState() {
		log.Info("winnig action")
		r.endOfGameWithWinner(obj.Key.Name)
		return true
	}

	if r.isTieState() {
		log.Info("tie state")
		r.endOfGameTie()
		return true
	}

	r.changeTurn()

	if err := r.setBoard(r.board); err != nil {
		r.log.WithError(err).Error("could not set the board")
		return true
	}

	return false
}

func (r *Referee) update(action types.Action) {
	for i := range r.board.Positions {
		pos := &r.board.Positions[i]
		if action.X == pos.X && action.Y == pos.Y {
			pos.Piece = r.board.Picked
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
		//Fatal
		r.board.Turn = ""
	}
}

func (r *Referee) validateAction(action types.Action) bool {
	if action.X < 1 || action.X > 16 || action.Y < 1 || action.Y > 16 {
		return false
	}

	for _, pos := range r.board.Positions {
		if pos.X == action.X && pos.Y == action.Y && pos.Piece != 0 {
			return false
		}

		if action.Picked == pos.Piece {
			return false
		}
	}

	if action.Picked < 1 || action.Picked > 16 {
		return false
	}

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

func (r *Referee) endOfGameWithWinner(winner string) {
	if r.p1.Name == winner {
		r.endOfGame(r.p1.Id, r.p2.Id, 3, 0)
	} else {
		r.endOfGame(r.p2.Id, r.p1.Id, 3, 0)
	}
}

func (r *Referee) endOfGameWithLoser(loser string) {
	if r.p1.Name == loser {
		r.endOfGame(r.p2.Id, r.p1.Id, 3, 0)
	} else {
		r.endOfGame(r.p1.Id, r.p2.Id, 3, 0)
	}
}

func (r *Referee) endOfGameTie() {
	r.endOfGame(r.p1.Id, r.p2.Id, 1, 1)

}

func (r *Referee) endOfGame(p1, p2 string, s1, s2 int) {

}

func (r *Referee) isWinState() bool {
	var l, s, c, h [11]int
	for _, pos := range r.board.Positions {
		if pos.Piece == 0 {
			continue
		}

		piece := r.board.Pieces[pos.Piece]

		if piece.Color == types.White {
			c[pos.X]++
			c[pos.Y+4]++
			if pos.X == pos.Y {
				c[9]++
			}
			if pos.X == 5-pos.Y {
				c[10]++
			}
		}

		if piece.Hole == types.Hollow {
			h[pos.X]++
			h[pos.Y+4]++
			if pos.X == pos.Y {
				h[9]++
			}
			if pos.X == 5-pos.Y {
				h[10]++
			}
		}

		if piece.Length == types.Short {
			l[pos.X]++
			l[pos.Y+4]++
			if pos.X == pos.Y {
				l[9]++
			}
			if pos.X == 5-pos.Y {
				l[10]++
			}
		}

		if piece.Shape == types.Round {
			s[pos.X]++
			s[pos.Y+4]++
			if pos.X == pos.Y {
				s[9]++
			}
			if pos.X == 5-pos.Y {
				s[10]++
			}
		}
	}

	for i := 1; i <= 10; i++ {
		if s[i]%4 == 0 || l[i]%4 == 0 || h[i]%4 == 0 || c[i]%4 == 0 {
			return true
		}
	}
	return false
}

func (r *Referee) isTieState() bool {
	for _, pos := range r.board.Positions {
		if pos.Piece == 0 {
			return false
		}
	}

	return true
}
