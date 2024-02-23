package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
    args := os.Args
    name := args[1]
    var tick_rate int64
    tick_rate, err := strconv.ParseInt(args[2], 10, 64);
    if err != nil {
        log.Panic(err)
	}
    fmt.Printf("fetus running with name: %v, tick_rate: %v\n", name, tick_rate)
    ticker := time.NewTicker(time.Duration(tick_rate * int64(time.Second)))
	defer ticker.Stop()
	done := make(chan bool)
	for {
		select {
		case <-done:
			fmt.Println("Done!")
			return
		case _ = <-ticker.C:
            fmt.Printf("%v: Tick!\n", name)
		}
	}
}
