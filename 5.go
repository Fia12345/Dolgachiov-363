package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"sync"
)

type Task struct {
	FilePath string
}

type Result struct {
	FilePath string
	Hash     string
	Err      error
}

func main() {
	files := []string{
		"file1.txt",
		"file2.txt",
		"file3.txt",
		"file4.txt",
		"file5.txt",
	}

	createTestFiles(files)

	workerCount := 2
	maxOpenFiles := 2

	tasks := make(chan Task, len(files))
	results := make(chan Result, len(files))
	fileSemaphore := make(chan struct{}, maxOpenFiles)

	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for task := range tasks {
				hash, err := calculateMD5(task.FilePath, fileSemaphore)
				results <- Result{task.FilePath, hash, err}
			}
		}(i)
	}

	for _, file := range files {
		tasks <- Task{file}
	}
	close(tasks)

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		if result.Err != nil {
			fmt.Printf("Error: %s - %v\n", result.FilePath, result.Err)
		} else {
			fmt.Printf("%s: %s\n", result.FilePath, result.Hash)
		}
	}

	cleanupTestFiles(files)
}

func calculateMD5(filePath string, fileSemaphore chan struct{}) (string, error) {
	fileSemaphore <- struct{}{}
	defer func() { <-fileSemaphore }()

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func createTestFiles(files []string) {
	content := []string{
		"42 братухо",
		"КИзяк s:=[] (блек файс)",
		"Астма",
		"Астана",
		"мудрые боги выдумали ноги",
	}

	for i, file := range files {
		os.WriteFile(file, []byte(content[i]), 0644)
	}
}

func cleanupTestFiles(files []string) {
	for _, file := range files {
		os.Remove(file)
	}
}