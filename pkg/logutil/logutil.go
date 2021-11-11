package logutil

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/labstack/gommon/color"
)

var _ io.Writer = PrettyPrinter{}

type PrettyPrinter struct {
	Enable bool
	Out    io.Writer
}

func (x PrettyPrinter) Write(in []byte) (n int, err error) {
	if !x.Enable {
		return os.Stdout.Write(in)
	}

	type LogContent struct {
		Message   string    `json:"message"`
		Level     string    `json:"level"`
		Timestamp time.Time `json:"time"`
	}

	// parse content
	content := LogContent{}
	if err = json.Unmarshal(in, &content); err != nil {
		outmsg := fmt.Sprintf("bad unmarshaling: %s", string(in))
		return os.Stdout.Write([]byte(outmsg))
	}

	// map tag color
	colorfuncmapper := map[string]func(msg interface{}, styles ...string) string{
		"debug": color.Cyan,
		"info":  color.Blue,
		"warn":  color.Yellow,
		"error": color.Red,
		"fatal": color.Red,
		"panic": color.Red,
	}

	if fn, ok := colorfuncmapper[strings.ToLower(content.Level)]; ok {
		content.Level = fn(content.Level)
	}

	// pretty print
	outmsg := fmt.Sprintf("%s %s %s %s",
		color.Green(content.Timestamp.Format(time.RFC3339)),
		content.Level,
		color.Yellow(strings.TrimSpace(content.Message)),
		color.Dim(string(in)),
	)

	return os.Stdout.Write([]byte(outmsg))
}
