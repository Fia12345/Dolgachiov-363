package main
import (
    "fmt"
	"net/http"
	"sync"
)

func checkURL(url string, wg *sync.WaitGroup, sem chan bool) {
	defer wg.Done()
	sem <- true
	defer func() {<-sem}()

	resp, err:= http.Get(url)
	if err != nil {
		fmt.Printf("Ошибка %s: %v\n", url, err )
		return
	}

	defer resp.Body.Close()

	fmt.Printf("%s: %d\n", url, resp.StatusCode)
}

	func main() {
		urls := []string{
			"https://google.com",
			"https://yandex.ru",
			"https://habr.com",
			"https://github.com",
			"https://steam.com",
			"https://Ktk45.ru",
		}
		 var wg sync.WaitGroup
		 sem := make(chan bool, 3)

		 for _, url := range urls{
			wg.Add(1)
			go checkURL(url, &wg, sem) 
		}
		wg.Wait()
	}

