package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	startTime  time.Time
	resTime    time.Duration
	urls       []string
	resultFile = "results.txt"
)

func main() {
	// читаем из файла список ресурсов
	resourcesFile, err := os.Open("resources.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	resourcesReader := bufio.NewReader(resourcesFile)
	for {
		line, _, err := resourcesReader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println(err)
				return
			}
		}
		urls = append(urls, string(line))
	}
	resourcesFile.Close()

	// создание нового файла results.txt
	file, err := os.Create(resultFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	// запись урлов в файл
	for _, url := range urls {
		fmt.Fprintln(file, url)
	}
	file.Close()

	// в бесконечном цикле последовательно отправляем GET-запросы на адреса из списка urls
	for {
		for i, url := range urls {
			startTime = time.Now()
			_, err := http.Get(url) // обрабатываем только ошибку
			if err != nil {
				fmt.Println(err)
				continue
			}
			resTime = time.Now().Sub(startTime) // вычисляем длительность получения ответа
			err = logWriter(resultFile, url, i, resTime.Nanoseconds())
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

// функция пишет результат (в миллисекундах) в файл result.txt
func logWriter(fileName string, url string, urlNum int, time int64) (err error) {
	value := "; " + strconv.Itoa(int(time/1000000)) + " ms"
	lines, err := read(fileName)
	if err != nil {
		return
	}
	lines[urlNum] += value
	err = write(lines, fileName)
	if err != nil {
		return
	}
	return
}

// чтение из файла
func read(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// запись в файл
func write(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}
