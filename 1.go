package main
import (
	"fmt"
    "sync"
	"time"
)
func main(){
var wg sync.WaitGroup
wg.Add(1)
go func(){
	defer wg.Done()
	for i:=1; i<=5;i++{
		fmt.Println(i)
		time.Sleep(time.Second)
	}
}()
wg.Wait()
fmt.Println("Конец")
}