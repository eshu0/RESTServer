package main

import (
	"fmt"

	ibot "github.com/eshu0/mybot/pkg/interfaces"
	"github.com/stianeikeland/go-rpio"
)

type TestBot struct {
	ibot.IMyBot
}

func NewTestBot(folder string) *TestBot {

	mbot := &TestBot{}
	fmt.Println("opening gpio")
	return mbot
}

func (bot *TestBot) Forwards() {
	fmt.Println("Forwards!")
}

func (bot *TestBot) Backwards() {
	fmt.Println("Backwards!")
}

func (bot *TestBot) Stop() {
	fmt.Println("Stop!")
}

func (bot *TestBot) SpinRight() {
	fmt.Println("Spinning Right!")
}

func (bot *TestBot) SpinLeft() {
	fmt.Println("Spinning Left!")
}

func (bot *TestBot) Close() {
	rpio.Close()
}

func (bot *TestBot) Hflip(b bool) {
}

func (bot *TestBot) Vflip(b bool) {
}

func (bot *TestBot) Capture() (string, error) {
	fmt.Println("Capturing Photo!")
	return "", nil
}
