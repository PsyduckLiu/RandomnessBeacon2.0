package watch

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
)

func WatchOutput(outputCh chan string, filename string) {
	watch, err := fsnotify.NewWatcher()
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from WatchOutput]Watch OUTPUT failed:%s", err))
	}
	defer watch.Close()

	err = watch.Add(filename)
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from WatchOutput]Watch OUTPUT failed:%s", err))
	}

	go func() {
		lastOutput := ""
		for {
			select {
			case ev := <-watch.Events:
				{
					if ev.Op&fsnotify.Write == fsnotify.Write {
						output := ReadFile(filename)
						if lastOutput != output {
							lastOutput = output
							outputCh <- lastOutput
							fmt.Println("===>[WatchOutput]New output is:", output)
						}
					}
				}
			case err := <-watch.Errors:
				{
					panic(fmt.Errorf("===>[ERROR from WatchOutput]Watch OUTPUT failed:%s", err))
				}
			}
		}
	}()

	select {}
}
