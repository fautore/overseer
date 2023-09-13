package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type process struct {
    command string
    arguments []string
    view *tview.TextView
}

func startChildProcess(outputTextView *tview.TextView, program string, arguments ...string) {
    cmd := exec.Command(program, arguments...)
    stdout, err := cmd.StdoutPipe()
    if err != nil { log.Fatal(err) }
    _, err = cmd.StdinPipe()
    if err != nil { log.Fatal(err) }
    _, err = cmd.StderrPipe()
    if err != nil { log.Fatal(err) }
    err = cmd.Start()
    if err != nil { log.Fatal(err) }
    go func () {
        reader := bufio.NewReader(stdout)
        for {
            line, err := reader.ReadString('\n')
            if err != nil {
                if err == io.EOF { break }
                fmt.Printf("Error reading process stdout: %v\n", err)
                break
            }
           fmt.Fprintf(outputTextView, "%s ", line) 
           outputTextView.ScrollToEnd()
        }
    }() 
}

func main() {
    // spawn application
    app := tview.NewApplication()
	inputField := tview.NewInputField().
		SetLabel("Enter Input (Press Enter to send): ").
		SetFieldWidth(40) 

    // start subprocesses
    processes := []process {
        { command: "fetus/fetus", arguments: []string{"1", "1"}, view: tview.NewTextView()}, 
        { command: "fetus/fetus", arguments: []string{"2", "3"}, view: tview.NewTextView()},
        { command: "fetus/fetus", arguments: []string{"3", "5"}, view: tview.NewTextView()},
        { command: "ping", arguments: []string{"www.google.com"}, view: tview.NewTextView()},
    }
    for _, p := range processes {
        if p.view != nil {
            p.view.SetChangedFunc(func() { app.Draw() })
            p.view.SetBorder(true)
            p.view.SetScrollable(true)
        } 
        go startChildProcess(p.view, p.command, p.arguments...)
    }

    inputField.SetDoneFunc(func(key tcell.Key) {
        if key == tcell.KeyEnter {
            inputText := inputField.GetText()
			inputField.SetText("")
            for _, p := range processes {
                p.view.SetText(inputText)
            }
        }
	})

	// Create a flex layout to split the windows
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(inputField, 0, 1, false)
    viewFlex := tview.NewFlex().SetDirection(tview.FlexRow)
    for _, p := range processes {
		viewFlex.AddItem(p.view, 0, 1, false)
    }
    flex.AddItem(viewFlex, 0, 10, false)

    if err := app.SetRoot(flex, true).SetFocus(inputField).Run(); err != nil {
	    fmt.Println(err)
	}
}
