package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

func readLinesChan(filePath string) chan []byte {
	resCh := make(chan []byte)

	go func() {
		f, err := os.Open(filePath)
		if err != nil {
			fmt.Println("Error:", err.Error())
			close(resCh)
			return
		}
		defer f.Close()

		reader := bufio.NewReader(f)
		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				if !errors.Is(err, io.EOF) {
					fmt.Println("Error:", err.Error())
				}

				if len(line) > 0 {
					resCh <- line
				}

				break
			}

			resCh <- line
		}

		close(resCh)
	}()

	return resCh
}
