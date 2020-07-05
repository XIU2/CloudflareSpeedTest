package main

import (
	"encoding/csv"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"log"
	"os"
	"sync"
)

func ExportCsv(filePath string, data [][]string) {
	fp, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("创建文件["+filePath+"]句柄失败,%v", err)
		return
	}
	defer fp.Close()
	w := csv.NewWriter(fp) //创建一个新的写入文件流
	w.WriteAll(data)
	w.Flush()
}

var pingTime int
var pingRoutine int
const ipEndWith uint8 = 1
type progressEvent int
const (
	NoAvailableIPFound progressEvent = iota
	AvailableIPFound
	NormalPing
)

func handleProgressGenerator(pb *pb.ProgressBar)func (e progressEvent){
	return func(e progressEvent) {
		switch e {
		case NoAvailableIPFound:
			pb.Add(pingTime)
		case AvailableIPFound:
			pb.Add(failTime)
		case NormalPing:
			pb.Increment()
		}
	}
}

func handleUserInput(){
	fmt.Println("请输入扫描协程数(数字越大越快,默认100):")
	fmt.Scanln(&pingRoutine)
	if pingRoutine<=0{
		pingRoutine=100
	}
	fmt.Println("请输入tcping次数(默认10):")
	fmt.Scanln(&pingTime)
	if pingTime<=0{
		pingTime=10
	}
}

func main(){
	handleUserInput()
	ips:=loadFirstIPOfRangeFromFile()
	pingCount:=len(ips)*pingTime
	bar := pb.StartNew(pingCount)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var data = make([][]string,0)
	data = append(data,[]string{"IP Address","Ping received","Ping time"})
	control := make(chan bool,pingRoutine)
	for _,ip :=range ips{
		wg.Add(1)
		control<-false
		handleProgress:=handleProgressGenerator(bar)
		go tcpingGoroutine(&wg,&mu,ip,pingTime, &data,control,handleProgress)
	}
	wg.Wait()
	bar.Finish()
	ExportCsv("./result.csv",data)
}
