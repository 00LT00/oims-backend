package gpu

import (
	"fmt"
	"github.com/NVIDIA/go-nvml/pkg/nvml"
	"log"
	"time"
)

var GPUs chan int

func init() {
	ret := nvml.Init()
	if ret != nvml.SUCCESS {
		log.Fatalf("Unable to initialize NVML: %v", nvml.ErrorString(ret))
	}
	GPUs = make(chan int)
	go updateEmptyGPU()
}

func updateEmptyGPU() {
	fmt.Println("update empty gpu")
	defer func() {
		ret := nvml.Shutdown()
		if ret != nvml.SUCCESS {
			log.Fatalf("Unable to shutdown NVML: %v", nvml.ErrorString(ret))
		}
	}()
	count, ret := nvml.DeviceGetCount()
	fmt.Println("gpu number: ", count)
	if ret != nvml.SUCCESS {
		log.Fatalf("Unable to get device count: %v", nvml.ErrorString(ret))
	}
	for {
		for i := 0; i < count; i++ {
			device, ret := nvml.DeviceGetHandleByIndex(i)
			if ret != nvml.SUCCESS {
				log.Fatalf("Unable to get device at index %d: %v", i, nvml.ErrorString(ret))
			}
			utilization, ret := device.GetUtilizationRates()
			if ret != nvml.SUCCESS {
				log.Fatalf("Unable to get Utilization of device at index %d: %v", i, nvml.ErrorString(ret))
			}
			fmt.Println("gpu: ", i, utilization)
			if utilization.Gpu <= 20 && utilization.Memory <= 30 {
				select {
				case GPUs <- i:
					fmt.Println("gpu", i, "was be acquired")
				default:
					fmt.Println("not use", i)
				}
			}
		}
		fmt.Println("try get empty")
		time.Sleep(time.Second * 20) //每3s更新一次空闲状态
	}
}
