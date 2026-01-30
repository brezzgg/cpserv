package log

import (
	"encoding/json"
	"fmt"
	"os"
)

type message struct {
	MessageType string `json:"type"`
	Message     any    `json:"msg"`
}

func Response(resp any) {
	b, err := json.Marshal(message{
		MessageType: "response",
		Message: resp,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	os.Exit(0)
}

func Error(err any) {
	b, err := json.Marshal(message{
		MessageType: "error",
		Message: err,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	os.Exit(1)
}
