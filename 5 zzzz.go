package main 
import (
    "fmt"
	"crypto/md5"
	"io"
	"os"
	"sync"
	"path/filepath"
)

type Task struct{
	FilePath string
}

type Result struct {
	FilePath string 
	Hash string 
	Err error
}

func main(){
	files:= []string{
		"file1.txt",
		"file2.txt",
		"file3.txt",
		"file4.txt",
	}
	workerCount:=3

	tasks:= make(chan Task, len(files))
	results:= make(chan Result, len(files))

	var wg sync.WaitGroup 
	for i:=0; i<workerCount;i++{
		wg.Add(1)
		go worker(i, tasks, results, &wg)
	}
	for _, file:=range files {
		tasks<-Task{FilePath:file}
	}
	close(tasks)

	go func(){
		wg.Wait()
		close(results)
	}()

	for result:=range results {
		if result.Err != nil {
			fmt.Printf("Ошибка: %s, %v \n", result.FilePath,result.Err)
		} else {
			fmt.Printf("%s, %v \n", result.FilePath,result.Hash)
		}
	}
}

func worker(id int, tasks<-chan Task, results chan<-Result, wg *sync.WaitGroup){
	defer wg.Done()
	for task:=range tasks{
		fmt.Printf("воркер %d обрабатывает: %s\n", id, task.FilePath)
		hash, err:=calculateMD5(task.FilePath)
		results<-Result{
			FilePath: task.FilePath,
			Hash: hash,
			Err: err,
		}
	}
}

func calculateMD5(filePath string)(string, error){
	file,err:=os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hash:=md5.New()
	if _, err:=io.Copy(hash,file);err!=nil{
		return "", err
	}
	return fmt.Sprintf("%x",hash.Sum(nil)), nil
}