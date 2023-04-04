package watch

import (
	"fmt"
	"generator/config"
	"time"
)

func WatchOutput(outputCh chan string, filename string) {
	config.DownloadFile("http://172.18.208.214/"+filename, "download/"+filename)
	lastOutput := config.ReadFile("download/" + filename)

	go func() {
		for {
			time.Sleep(1 * time.Second)

			config.DownloadFile("http://172.18.208.214/"+filename, "download/"+filename)
			output := config.ReadFile("download/" + filename)

			if lastOutput != output && output != "" {
				lastOutput = output
				outputCh <- output
				fmt.Println("===>[WatchOutput]New output is:", output)
			}
		}
	}()
}
