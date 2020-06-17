package referee

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/Gimulator-Games/quarto-judge/types"
	"github.com/Gimulator/client-go"
)

const (
	actionType   = "action"
	verdictType  = "verdict"
	registerType = "register"
	timerType    = "timer"
	endGameType  = "end-of-game"

	namespace   = "quarto"
	name        = "referee"
	boardName   = "board"
	apiTimeWait = time.Second * 1
)

type cli struct {
	roomID string
	*client.Client
}

func newCli(ch chan client.Object, roomID string) (cli, error) {
	c, err := client.NewClient(ch)
	if err != nil {
		return cli{}, err
	}

	if err := c.Watch(client.Key{
		Type:      actionType,
		Namespace: namespace,
		Name:      "",
	}); err != nil {
		return cli{}, err
	}

	return cli{
		Client: c,
		roomID: roomID,
	}, nil
}

func (c *cli) receiptPlayers(ctx context.Context) (*client.Object, *client.Object, int) {
	ticker := time.NewTicker(time.Second * 2)
	done := ctx.Done()

	key := client.Key{
		Namespace: namespace,
		Type:      registerType,
		Name:      "",
	}

	for {
		select {
		case <-ticker.C:
			if o1, o2, n := c.findPlayers(key); n == 2 {
				return o1, o2, n
			}
		case <-done:
			if o1, o2, n := c.findPlayers(key); n == 2 {
				return o1, o2, n
			} else if n == 1 {
				c.setEndOfGameWithOnePlayer(o1)
				return o1, nil, n
			} else {
				c.setEndOfGameWithoutPlayer()
				return nil, nil, n
			}
		}
	}
}

func (c *cli) findPlayers(key client.Key) (*client.Object, *client.Object, int) {
	res, err := c.Find(key)
	if err != nil {
		return nil, nil, 0
	}

	switch len(res) {
	case 0:
		return nil, nil, 0
	case 1:
		return &res[0], nil, 1
	case 2:
		return &res[0], &res[1], 2
	default:
		return nil, nil, 0
	}
}

func (c *cli) setEndOfGameWithoutPlayer() error {
	res := result{
		RoomID:  c.roomID,
		Status:  "FAIL",
		Message: "no players registred",
	}
	return c.setEndOfGame(res)
}

func (c *cli) setEndOfGameWithOnePlayer(obj *client.Object) error {
	scores := make(map[string]map[int]int)
	scores[obj.Meta.Owner] = map[int]int{1: 1}

	res := result{
		RoomID:  c.roomID,
		Status:  "SUCCESS",
		Message: "",
		Scores:  scores,
	}
	return c.setEndOfGame(res)
}

func (c *cli) setEndOfGame(res result) error {
	bytes, err := json.Marshal(res)
	if err != nil {
		return err
	}

	key := client.Key{
		Type:      endGameType,
		Namespace: namespace,
		Name:      name,
	}

	for {
		err := c.Set(key, string(bytes))
		if err == nil {
			os.Exit(0)
			return nil
		}

		time.Sleep(apiTimeWait)
	}
}

func (c *cli) setBoard(b types.Board) error {
	val, err := json.Marshal(b)
	if err != nil {
		return err
	}

	key := client.Key{
		Namespace: namespace,
		Type:      verdictType,
		Name:      boardName,
	}

	for {
		if err := c.Set(key, string(val)); err == nil {
			return nil
		}
		time.Sleep(apiTimeWait)
	}
}

type result struct {
	RoomID  string                 `json:"run_id"`
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Scores  map[string]map[int]int `json:"scores"`
}
