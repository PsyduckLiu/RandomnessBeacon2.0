package config

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
)

func DownloadFile(inUrl string, outFile string) {
	// Get the data
	resp, err := http.Get(inUrl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	out, err := os.Create(outFile)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}
}

func CopyFile(srcFile, destFile string) (int64, error) {
	file1, err := os.Open(srcFile)
	if err != nil {
		return 0, err
	}
	file2, err := os.OpenFile(destFile, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return 0, err
	}
	defer file1.Close()
	defer file2.Close()
	return io.Copy(file2, file1)
}

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
