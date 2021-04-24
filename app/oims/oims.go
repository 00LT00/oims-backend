package oims

import (
	"oims/service"
	"oims/service/gpu"
	"os/exec"
	"strconv"
)

var conf = service.Conf
var jobs chan string

func run(id string) {
	cmd := exec.Command("/usr/bin/python3 /oims/model_inspect.py",
		"--output_image_dir", conf.Path.Result,
		"--input_image", conf.Path.History+"/"+id+"/*.jpg",
		"--gpu", strconv.Itoa(<-gpu.GPUs),
	)
	err:=cmd.Wait()
	if err != nil {
		jobs<- id
	}
}

func Add(id string){
	jobs<- id
}

func init() {
	jobs = make(chan string)
	go func() {
		for {
			select {
			case id := <- jobs:
					go run(id)

			}
		}
	}()
}
