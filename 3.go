package main
import (
 "fmt"
 "time"
)

func main(){
tick:= time.Tick(200*time.Millisecond)
requests:=make(chan int,15)
for i:=1;i<=15;i++{
	requests <-i
}
close(requests)
for req:=range requests{
	<-tick
	fmt.Printf("Запрос: %d в %v\n",req,time.Now())
}
}
