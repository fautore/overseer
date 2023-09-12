package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type process struct {
    command string
    arguments []string
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
        }
    }() 
}

func main() {
    // spawn application
    app := tview.NewApplication()
	inputField := tview.NewInputField().
		SetLabel("Enter Input (Press Enter to send): ").
		SetFieldWidth(40).
		SetAcceptanceFunc(tview.InputFieldMaxLength(100))

    inputField.SetDoneFunc(func(key tcell.Key) {
        if key == tcell.KeyEnter {
            inputText := inputField.GetText()
			os.Stdin.WriteString(inputText)
			inputField.SetText("")
        }
	})

	// Create a text view for subprocess output
	outputTextView := tview.NewTextView().
        SetChangedFunc(func() {
		    app.Draw()
	})
    outputTextView.SetBorder(true)

    // start subprocesses
    processes := []process {
        { command: "fetus/fetus", arguments: []string{"1", "1"}}, 
        { command: "fetus/fetus", arguments: []string{"2", "3"}},
        { command: "fetus/fetus", arguments: []string{"3", "5"}},
    }
    for _, p := range processes {
        go startChildProcess(outputTextView, p.command, p.arguments...)
    }

	// Create a flex layout to split the windows
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(inputField, 0, 1, false).
		AddItem(outputTextView, 0, 2, false)

	if err := app.SetRoot(flex, true).Run(); err != nil {
		fmt.Println(err)
	}


    
}
