package main

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/Gimulator-Games/quarto-judge/referee"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)

	formatter := &logrus.TextFormatter{
		TimestampFormat:  "2006-01-02 15:04:05",
		FullTimestamp:    true,
		PadLevelText:     true,
		QuoteEmptyFields: true,
		ForceQuote:       false,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return "", fmt.Sprintf(" %s:%d\t", path.Base(f.File), f.Line)
		},
	}
	logrus.SetFormatter(formatter)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	roomID := os.Getenv("ROOM_ID")
	if roomID == "" {
		panic("set 'ROOM_ID' environment variable to send result of game")
	}

	r, err := referee.NewReferee(roomID)
	if err != nil {
		panic(err)
	}

	if err := r.Start(); err != nil {
		panic(r)
	}
}
