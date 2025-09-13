package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	sources := []string{"DB1", "DB2", "DB3", "API1", "API2"}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result := parallelSearch(ctx, sources, "test query")

	if result != "" {
		fmt.Printf("Найден результат: %s\n", result)
	} else {
		fmt.Println("Результат не найден")
	}
}

func parallelSearch(ctx context.Context, sources []string, query string) string {
	resultChan := make(chan string, len(sources))
	var wg sync.WaitGroup

	for _, source := range sources {
		wg.Add(1)
		go func(src string) {
			defer wg.Done()
			if res, err := searchInSource(ctx, src, query); err == nil {
				select {
				case resultChan <- res:
				default:
				}
			}
		}(source)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	select {
	case res := <-resultChan:
		return res
	case <-ctx.Done():
		return ""
	}
}

func searchInSource(ctx context.Context, source, query string) (string, error) {
	delay := time.Duration(rand.Intn(2000)+500) * time.Millisecond

	select {
	case <-time.After(delay):
		return fmt.Sprintf("%s: результат для '%s'", source, query), nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}