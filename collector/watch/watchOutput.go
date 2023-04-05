package watch

import (
	"collector/config"
	"fmt"
	"time"
)

func WatchOutput(outputCh chan string, filename string) {
	boardIP := config.GetBoardIP()

	config.DownloadFile("http://"+boardIP+"/"+filename, "download/"+filename)
	lastOutput := config.ReadFile("download/" + filename)

	go func() {
		for {
			time.Sleep(1 * time.Second)

			config.DownloadFile("http://"+boardIP+"/"+filename, "download/"+filename)
			output := config.ReadFile("download/" + filename)

			if lastOutput != output && output != "" {
				lastOutput = output
				outputCh <- output
				fmt.Println("===>[WatchOutput]New output is:", output)
			}
		}
	}()
}
