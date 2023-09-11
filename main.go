package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
)

type process struct {
    command string
    arguments []string
}

func readStuff(scanner *bufio.Scanner) {
    for scanner.Scan() {
        fmt.Println(scanner.Text())
    }
    if err := scanner.Err(); err != nil {
        fmt.Fprintln(os.Stderr, "reading standard input:", err)
    }
}

func startChildProcess(program string, arguments ...string) {
    cmd := exec.Command(program, arguments...)
    out, err := cmd.StdoutPipe()
    err = cmd.Start()
    if err != nil {
        log.Fatal(err)
    }
    scanner := bufio.NewScanner(out)
    defer cmd.Wait()
    go readStuff(scanner)
}

func main() {
    processes := []process {
        { command: "fetus/fetus", arguments: []string{"1", "1"}}, 
        { command: "fetus/fetus", arguments: []string{"2", "3"}},
        { command: "fetus/fetus", arguments: []string{"3", "5"}},
    }
    for _, p := range processes {
        go startChildProcess(p.command, p.arguments...)
    }
    for {
    }
    //startChildProcess(exec.Command("fetus/fetus", "half", "0.5")) 
}
