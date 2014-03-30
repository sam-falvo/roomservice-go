package main

import (
	"github.com/nsf/termbox-go"
	"github.com/sam-falvo/soti/console"
	"github.com/sam-falvo/soti/window"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type Message struct {
	Name, Description  string
	Globals, Procedure string
}

type MessageBase struct {
	msgs   []*Message
	tagged map[*Message]bool
}

func NewBase() *MessageBase {
	mb := new(MessageBase)
	mb.tagged = make(map[*Message]bool)
	return mb
}

func drawListMenu(listTitle, appTitle string, mb *MessageBase) {
	w, h := termbox.Size()
	wnd := window.New(0, 0, w, h)
	wnd.SetFgPen(termbox.ColorWhite)
	wnd.SetBgPen(termbox.ColorBlack)
	wnd.Clear()

	console.PrintMenuTitle(wnd, listTitle, appTitle)
	window.AtPrint(wnd, 2, 2, "C  Create   U  Untag All     E  Edit Tagged")
	window.AtPrint(wnd, 2, 3, "X  Exit     T  Tag/Untag     D  Delete Tagged")

	wnd.MoveTo(2, 4)
	wnd.BoxTo(w-2, h)

	wnd2 := window.New(3, 5, w-3, h-1)
	wnd2.SetFgPen(termbox.ColorWhite)
	wnd2.SetBgPen(termbox.ColorBlack)
	wnd2.Clear()

	for y, m := range mb.msgs {
		checkmark := " "
		fg := termbox.ColorWhite
		if mb.tagged[m] {
			fg |= termbox.AttrBold
			checkmark = "/"
		}
		wnd2.SetFgPen(fg)
		window.AtPrint(wnd2, 0, y, checkmark)
		window.AtPrint(wnd2, 2, y, m.Name)
		window.AtPrint(wnd2, 11, y, m.Description)
		y++
	}

	termbox.Flush()
}

func deleteTagged(mb *MessageBase) {
	ms := make([]*Message, 0)
	for _, m := range mb.msgs {
		if !mb.tagged[m] {
			ms = append(ms, m)
		}
	}
	mb.msgs = ms
}

func messageList(title, appTitle string, mb *MessageBase) {
	done := false
	for !done {
		drawListMenu(title+" List", appTitle, mb)
		sel := console.WaitAscii()
		switch sel {
		case 'X', 'x':
			done = true
		case 'C', 'c':
			createMessage(title, appTitle, mb)
		case 'U', 'u':
			mb.tagged = make(map[*Message]bool)
		case 'T', 't':
			tagMessages(mb)
		case 'D', 'd':
			deleteTagged(mb)
		case 'E', 'e':
			for m, _ := range mb.tagged {
				editMessage(title, appTitle, m)
			}
		}
	}
}

func drawTagMenu(mb *MessageBase, tagged map[*Message]bool, cursor int) {
	w, h := termbox.Size()
	wnd := window.New(4, 2, w-4, h-2)
	wnd.SetFgPen(termbox.ColorWhite)
	wnd.SetBgPen(termbox.ColorBlack)
	wnd.Clear()
	wnd.MoveTo(0, 0)
	wnd.BoxTo(w-8, h-4)

	window.AtPrint(wnd, 1, 1, "X=Cancel   W=OK   Spc=Tag/Untag   J,K=Down,Up")
	for i, m := range mb.msgs {
		fg := termbox.ColorWhite
		bg := termbox.ColorBlack
		checkmark := " "
		if i == cursor {
			fg, bg = bg, fg
		}
		if mb.tagged[m] {
			fg |= termbox.AttrBold
			checkmark = "/"
		}
		wnd.SetFgPen(fg)
		wnd.SetBgPen(bg)
		window.AtPrint(wnd, 1, 3+i, checkmark)
		window.AtPrint(wnd, 3, 3+i, m.Name)
		window.AtPrint(wnd, 13, 3+i, m.Description)
		i++
	}
	termbox.Flush()
}

func currentlyTagged(mb *MessageBase) map[*Message]bool {
	tagged := make(map[*Message]bool)
	for _, m := range mb.msgs {
		if mb.tagged[m] {
			tagged[m] = true
		}
	}
	return tagged
}

func tagMessages(mb *MessageBase) {
	tagged := currentlyTagged(mb)
	done := false
	cursor := 0
	m := len(mb.msgs)

	for !done {
		drawTagMenu(mb, tagged, cursor)
		sel := console.WaitAscii()
		switch sel {
		case 'X', 'x':
			done = true
		case 'W', 'w':
			done = true
			mb.tagged = tagged
		case 'J', 'j':
			if (cursor + 1) < m {
				cursor++
			}
		case 'K', 'k':
			if cursor > 0 {
				cursor--
			}
		case ' ':
			s := mb.msgs[cursor]
			tagged[s] = !tagged[s]
		}
	}
}

func drawEditMenu(title, appTitle string, m *Message) {
	w, h := termbox.Size()
	wnd := window.New(0, 0, w, h)
	wnd.SetFgPen(termbox.ColorWhite)
	wnd.SetBgPen(termbox.ColorBlack)
	wnd.Clear()

	console.PrintMenuTitle(wnd, title, appTitle)
	window.AtPrint(wnd, 2, 2, "N  Name")
	console.PrintField(wnd, w-10, w-2, m.Name)

	window.AtPrint(wnd, 2, 3, "D  Description ")
	console.PrintField(wnd, w-60, w-2, m.Description)

	window.AtPrint(wnd, 2, 4, "G  Globals (opens editor)")
	window.AtPrint(wnd, 2, 5, "P  Procedure (opens editor)")

	window.AtPrint(wnd, 2, 7, "W  Exit and save changes")
	window.AtPrint(wnd, 2, 8, "X  Exit and discard chanegs")

	termbox.Flush()
}

func editMessage(title, appTitle string, m *Message) (commit bool) {
	t := new(Message)
	*t = *m
	done := false
	for !done {
		drawEditMenu(title, appTitle, t)
		sel := console.WaitAscii()
		switch sel {
		case 'N', 'n':
			t.Name = editName(t.Name)
		case 'D', 'd':
			t.Description = editDesc(t.Description)
		case 'X', 'x':
			done = true
		case 'W', 'w':
			if strings.TrimSpace(t.Name) != "" {
				done = true
				*m = *t
				commit = true
			}
		case 'G', 'g':
			t.Globals = editTextArea(t.Globals)
		case 'P', 'p':
			t.Procedure = editTextArea(t.Procedure)
		}
	}
	return
}

func createMessage(title, appTitle string, mb *MessageBase) {
	m := new(Message)
	if editMessage(title+" Creation", appTitle, m) {
		mb.msgs = append(mb.msgs, m)
	}
}

func editName(oldName string) string {
	w, _ := termbox.Size()
	wnd := window.New(w-10, 2, w-2, 3)
	fd := &console.FieldDesc{
		Window: wnd,
		Value:  oldName,
	}
	console.EditField(fd)
	return fd.Value
}

func editDesc(oldDesc string) string {
	w, _ := termbox.Size()
	wnd := window.New(w-60, 3, w-2, 4)
	fd := &console.FieldDesc{
		Window:      wnd,
		Value:       oldDesc,
		AllowSpaces: true,
	}
	console.EditField(fd)
	return fd.Value
}

func editTextArea(oldText string) string {
	tmpFilename := ""
	defer func() {
		if tmpFilename != "" {
			os.Remove(tmpFilename)
		}
	}()

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = os.Getenv("VISUAL")
	}
	if editor == "" {
		return oldText
	}

	tmpf, err := ioutil.TempFile("", "roomservice-go-")
	if err != nil {
		return oldText
	}
	tmpFilename = tmpf.Name()

	_, err = tmpf.Write([]byte(oldText))
	tmpf.Close()
	if err != nil {
		return oldText
	}

	termbox.Close()
	cmd := exec.Command(editor, tmpFilename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	termbox.Init()
	if err != nil {
		return oldText
	}

	newText, err := ioutil.ReadFile(tmpFilename)
	if err != nil {
		return oldText
	}

	return string(newText)
}
