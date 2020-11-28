package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	Request "github.com/eshu0/RESTServer/pkg/request"
	ibot "github.com/eshu0/mybot/pkg/interfaces"
)

type Rbot struct {
	mbot ibot.IMyBot
}

func NewRbot(folder string) Rbot {
	mbot := NewTestBot(folder)
	rb := Rbot{mbot: mbot}
	return rb
}

func (bot Rbot) Stop(request Request.ServerRequest) {
	// stupid CORS - Means that for dev this is a pain
	// should be removed for PROD or for internet access
	request.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Println("Rbot - Stop")
	bot.mbot.Stop()
}

func (bot Rbot) Forwards(request Request.ServerRequest) {
	// stupid CORS - Means that for dev this is a pain
	// should be removed for PROD or for internet access
	request.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Println("Rbot - Forwards")
	bot.mbot.Forwards()
}

func (bot Rbot) Backwards(request Request.ServerRequest) {
	// stupid CORS - Means that for dev this is a pain
	// should be removed for PROD or for internet access
	request.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Println("Rbot - Backwards")
	bot.mbot.Backwards()
}

func (bot Rbot) TurnLeft(request Request.ServerRequest) {
	request.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	fmt.Println("Rbot - TurnLeft")
	bot.mbot.SpinLeft()
	time.Sleep(200 * time.Millisecond)
	bot.mbot.Stop()

}

func (bot Rbot) TurnRight(request Request.ServerRequest) {
	// stupid CORS - Means that for dev this is a pain
	// should be removed for PROD or for internet access
	request.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Println("Rbot - TurnRight")
	bot.mbot.SpinRight()
	time.Sleep(200 * time.Millisecond)
	bot.mbot.Stop()
}

func (bot Rbot) SpinLeft(request Request.ServerRequest) {
	// stupid CORS - Means that for dev this is a pain
	// should be removed for PROD or for internet access
	request.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	fmt.Println("Rbot - SpinLeft")
	bot.mbot.SpinLeft()
}

func (bot Rbot) SpinRight(request Request.ServerRequest) {
	// stupid CORS - Means that for dev this is a pain
	// should be removed for PROD or for internet access
	request.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Println("Rbot - SpinRight")
	bot.mbot.SpinRight()
}

func (bot Rbot) Capture(request Request.ServerRequest) {
	// stupid CORS - Means that for dev this is a pain
	// should be removed for PROD or for internet access
	request.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Println("Rbot - Capture")
	s, err := bot.mbot.Capture()

	if err == nil {
		fmt.Println("Rbot - Capture Success")
		fmt.Println(s)
	} else {
		fmt.Println("Rbot - Error:")
		fmt.Println(err.Error())
	}
}

func (bot Rbot) FlipHR(request Request.ServerRequest) {
	// stupid CORS - Means that for dev this is a pain
	// should be removed for PROD or for internet access
	request.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Println("Rbot - FlipHR")
	bot.mbot.Hflip(false)
}

func (bot Rbot) FlipH(request Request.ServerRequest) {
	// stupid CORS - Means that for dev this is a pain
	// should be removed for PROD or for internet access
	request.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Println("Rbot - FlipH")
	bot.mbot.Hflip(true)
}

func (bot Rbot) FlipVR(request Request.ServerRequest) {
	// stupid CORS - Means that for dev this is a pain
	// should be removed for PROD or for internet access
	request.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Println("Rbot - FlipVR")
	bot.mbot.Vflip(false)
}

func (bot Rbot) FlipV(request Request.ServerRequest) {
	// stupid CORS - Means that for dev this is a pain
	// should be removed for PROD or for internet access
	request.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Println("Rbot - FlipV")
	bot.mbot.Vflip(true)
}

type Data struct {
	Total    int
	Pos      int
	FileName string
	Alt      string
}

func (bot Rbot) Command(request Request.ServerRequest) {
	fmt.Println("Rbot - Command")
	request.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	//var files []string

	root := "./images/"
	count := 1
	res := []*Data{}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		//name := filepath.Ext("main.test.js")
		base := filepath.Base(path)

		if !info.IsDir() && len(base) > 1 && base[0] != '.' {
			count++
			i := Data{FileName: path, Pos: count, Alt: info.Name()} //filepath.Base(path)}
			res = append(res, &i)
		}

		return nil
	})
	if err != nil {
		panic(err)
	}

	for _, i := range res {
		fmt.Println(i.FileName)
		i.Total = len(res)
	}
	err = request.Template.Execute(request.Writer, res)
	if err != nil {
		fmt.Printf("Rbot - Command Error : %s\n", err.Error())
		return
	}
}

func (bot Rbot) Close() {
	if bot.mbot != nil {
		bot.mbot.Stop()
		bot.mbot.Close()
	}
}
