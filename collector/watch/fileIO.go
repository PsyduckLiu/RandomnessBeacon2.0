package watch

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func ReadFile(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from ReadFile]Read file failed:%s", err))
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	line, _, _ := reader.ReadLine()

	return string(line)
}

func WriteFile(filename string, data string) {
	time.Sleep(5 * time.Second)

	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from WriteFile]Write file failed:%s", err))
	}
	n, err := f.Write([]byte(data))
	if err != nil && n < len(data) {
		panic(fmt.Errorf("===>[ERROR from WriteFile]Write file failed:%s", err))
	}
	err = f.Close()
	if err != nil {
		panic(fmt.Errorf("===>[ERROR from WriteFile]Write file failed:%s", err))
	}

	fmt.Println("==>[WriteFile] Write ", data)
}
