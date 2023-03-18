package main

import (
	"RB/watch"
	"fmt"
)

func main() {
	outputCh := make(chan string)
	go watch.WatchOutput(outputCh, "../output")

	for {
		select {
		case newOutput := <-outputCh:
			{
				fmt.Println("===>[Collector]New output is:", newOutput)
			}
		}
	}
}
