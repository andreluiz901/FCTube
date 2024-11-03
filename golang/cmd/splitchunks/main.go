package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

const chunkSize = 1 * 1024 * 1024

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <mp4-file-path>")
		return
	}

	mp4Path := os.Args[1]

	file, err := os.Open(mp4Path)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	chunkDir := filepath.Dir(mp4Path)
	fmt.Printf("Chunk directory: %s\n", chunkDir)

	chunkCount := 0

	for {
		chunkFileName := filepath.Join(chunkDir, strconv.Itoa(chunkCount)+".chunk")

		chunkFile, err := os.Create(chunkFileName)
		if err != nil {
			fmt.Printf("Error creating chunk file: %v\n", err)
			return
		}

		_, err = io.CopyN(chunkFile, file, chunkSize)
		if err != nil {
			if err == io.EOF {
				chunkFile.Close()
				fmt.Println("Finished splitting the file into chunks.")
				break
			} else {
				fmt.Printf("Error copying chunk: %v\n", err)
				chunkFile.Close()
				return
			}
		}

		chunkFile.Close()
		fmt.Printf("Created chunk: %s\n", chunkFileName)

		chunkCount++
	}
}
