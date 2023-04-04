package watch

import (
	"RB/config"
	"fmt"
	"time"
)

func WatchOutput(filename string) {
	lastOutput := ""

	go func() {
		for {
			time.Sleep(1 * time.Second)

			config.DownloadFile("http://172.18.208.214/"+filename, "download/"+filename)
			output := config.ReadFile("download/" + filename)

			if lastOutput != output && output != "" {
				lastOutput = output
				// outputCh <- output
				fmt.Println("===>[WatchOutput]New output is:", output)
			}
		}
	}()
}

// watch for local file

// func WatchOutput(outputCh chan string, filename string) {
// 	watch, err := fsnotify.NewWatcher()
// 	if err != nil {
// 		panic(fmt.Errorf("===>[ERROR from WatchOutput]Watch OUTPUT failed:%s", err))
// 	}
// 	defer watch.Close()

// 	err = watch.Add(filename)
// 	if err != nil {
// 		panic(fmt.Errorf("===>[ERROR from WatchOutput]Watch OUTPUT failed:%s", err))
// 	}

// 	go func() {
// 		lastOutput := ""
// 		for {
// 			select {
// 			case ev := <-watch.Events:
// 				{
// 					if ev.Op&fsnotify.Write == fsnotify.Write {
// 						output := ReadFile(filename)
// 						if lastOutput != output && output != "" {
// 							lastOutput = output
// 							outputCh <- output
// 							fmt.Println("===>[WatchOutput]New output is:", output)
// 						}
// 					}
// 				}
// 			case err := <-watch.Errors:
// 				{
// 					panic(fmt.Errorf("===>[ERROR from WatchOutput]Watch OUTPUT failed:%s", err))
// 				}
// 			}
// 		}
// 	}()

// 	select {}
// }
