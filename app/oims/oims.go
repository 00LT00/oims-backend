package oims

import (
	"oims/service"
	"oims/service/gpu"
	"os/exec"
	"strconv"
)

var conf = service.Conf
var jobs chan string
var serv = service.Service

func run(id string) {
	cmd := exec.Command("/usr/bin/python3 /oims/model_inspect.py",
		"--output_image_dir", conf.Path.Result,
		"--input_image", conf.Path.History+"/"+id+"/*.jpg",
		"--gpu", strconv.Itoa(<-gpu.GPUs),
	)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	err := cmd.Start()
	if err != nil {
		panic(err)
	}

	go func(id string) {
		serv.Logger.Printf("-----------------------  %s  -----------------------\n", id)
		buf := make([]byte, 1024)
		temp := make([]byte, 1024)
		for {
			_, err := stdout.Read(temp)
			buf = append(buf, temp...)
			if err != nil {
				break
			}
		}
		serv.Logger.Println(buf)
		serv.Logger.Printf("-----------------------  %s  -----------------------\n\n", id)
	}(id)

	go func(id string) {
		serv.ErrLogger.Printf("-----------------------  %s  -----------------------\n", id)
		buf := make([]byte, 1024)
		temp := make([]byte, 1024)
		for {
			_, err := stderr.Read(temp)
			buf = append(buf, temp...)
			if err != nil {
				break
			}
		}
		serv.ErrLogger.Println(buf)
		serv.ErrLogger.Printf("-----------------------  %s  -----------------------\n\n", id)
	}(id)
	err = cmd.Wait()
	if err != nil {
		serv.ErrLogger.Println(id, " Remeasuring...")
		jobs <- id
	}
}

func Add(id string) {
	jobs <- id
}

func init() {
	jobs = make(chan string)
	go func() {
		for {
			select {
			case id := <-jobs:
				go run(id)

			}
		}
	}()
}
