package main

import (
	"context"
	"log"
	"machine"

	keyboard "github.com/sago35/tinygo-keyboard"
	"github.com/sago35/tinygo-keyboard/keycodes/jp"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	d := keyboard.New([]machine.Pin{
		machine.D1,
		machine.D2,
	}, []machine.Pin{
		machine.D8,
		machine.D9,
		machine.D10,
	}, [][][]keyboard.Keycode{
		{
			{jp.KeyT, jp.KeyY},
			{jp.KeyI, jp.KeyG},
			{jp.KeyN, jp.KeyO},
		},
	})

	d.Debug = true
	return d.Loop(context.Background())
}
