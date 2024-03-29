package oims

import (
	"fmt"
	"oims/service"
	"oims/service/gpu"
	"os/exec"
	"strconv"
)

var conf = service.Conf
var jobs chan string
var serv = service.Service
var cancelMap map[string]int

func run(id string) {
	fmt.Println("init run", id)
	gpuID := strconv.Itoa(<-gpu.GPUs)
	fmt.Println(id, " get gpu:", gpuID)
	cmd := exec.Command("/usr/bin/python3", "/oims/model_inspect.py",
		"--output_image_dir", conf.Path.Result+"/",
		"--input_image", conf.Path.History+"/"+id+"/*.jpg",
		"--gpu", gpuID,
	)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	err := cmd.Start()
	if err != nil {
		panic(err)
	}

	go func(id string) {
		buf := make([]byte, 1024)
		temp := make([]byte, 1024)
		buf = append(buf, []byte(fmt.Sprintf("\n-----------------------  %s  -----------------------\n", id))...)
		for {
			_, err := stdout.Read(temp)
			buf = append(buf, temp...)
			if err != nil {
				break
			}
		}
		buf = append(buf, []byte(fmt.Sprintf("\n-----------------------  %s  -----------------------\n", id))...)
		serv.Logger.Println(string(buf))
	}(id)

	go func(id string) {
		buf := make([]byte, 1024)
		temp := make([]byte, 1024)
		buf = append(buf, []byte(fmt.Sprintf("\n-----------------------  %s  -----------------------\n", id))...)
		for {
			_, err := stderr.Read(temp)
			buf = append(buf, temp...)
			if err != nil {
				break
			}
		}
		buf = append(buf, []byte(fmt.Sprintf("\n-----------------------  %s  -----------------------\n", id))...)
		serv.ErrLogger.Println(string(buf))
	}(id)
	err = cmd.Wait()
	if err != nil {
		serv.ErrLogger.Println(" Remeasuring ", id)
		Add(id)
	}
}

func Add(id string) {
	if cancelMap[id] >= 5 {
		serv.Logger.Println(id, "was cancel")
		fmt.Println(id, "was cancel")
		delete(cancelMap, id)
		return
	}
	fmt.Println("add", id)
	jobs <- id
}

func Cancel(id string) {
	cancelMap[id] = 5
}

func init() {
	jobs = make(chan string)
	cancelMap = make(map[string]int)
	go func() {
		for {
			select {
			case id := <-jobs:
				fmt.Println("run", id)
				go run(id)
			}
		}
	}()
}
