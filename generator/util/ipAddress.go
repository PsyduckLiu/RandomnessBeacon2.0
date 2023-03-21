package util

import (
	"bufio"
	"fmt"
	"os"
)

func GetIPAddress(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from ReadFile]Read file failed:%s", err))
	}
	defer file.Close()

	var ipList []string
	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		ipList = append(ipList, string(line))
	}

	return ipList
}
