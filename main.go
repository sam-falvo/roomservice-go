package main

import (
	"github.com/nsf/termbox-go"
	"github.com/sam-falvo/soti/console"
	"github.com/sam-falvo/soti/window"
)

const (
	kTITLE    = "RoomService"
	kBY       = "Samuel A. Falvo II"
	kVERSION  = "Release 1.0"
	kAppTitle = "RoomService V1.0"
)

var (
	stimuli   *MessageBase
	responses *MessageBase
)

func drawMainMenu() {
	w, h := termbox.Size()
	wnd := window.New(0, 0, w, h)
	wnd.SetFgPen(termbox.ColorWhite)
	wnd.SetBgPen(termbox.ColorBlack)
	wnd.Clear()

	y := h / 5
	window.AtPrint(wnd, (w/2)-len(kTITLE)/2, y, kTITLE)
	window.AtPrint(wnd, (w/2)-len(kBY)/2, y+1, kBY)
	window.AtPrint(wnd, (w/2)-len(kVERSION)/2, y+2, kVERSION)

	y = h / 2
	window.AtPrint(wnd, 2, y, "S  List / Edit Stimuli")
	window.AtPrint(wnd, 2, y+1, "R  List / Edit Responses")
	window.AtPrint(wnd, 2, y+2, "X  Exit")
	termbox.Flush()
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	stimuli = NewBase()
	responses = NewBase()

	done := false
	for !done {
		drawMainMenu()
		sel := console.WaitAscii()
		switch sel {
		case 'X', 'x':
			done = true
		case 'S', 's':
			messageList("Stimulus", kAppTitle, stimuli)
		case 'R', 'r':
			messageList("Response", kAppTitle, responses)
		}
	}
}
