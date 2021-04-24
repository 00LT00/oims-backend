package gpu

import (
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
	go updateEmptyGPU()
}

func updateEmptyGPU() {
	defer func() {
		ret := nvml.Shutdown()
		if ret != nvml.SUCCESS {
			log.Fatalf("Unable to shutdown NVML: %v", nvml.ErrorString(ret))
		}
	}()
	count, ret := nvml.DeviceGetCount()
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
			if utilization.Gpu <= 20 {
				select {
				case GPUs<-i:
				default:
				}
			}
		}
		time.Sleep(time.Second * 3) //每3s更新一次空闲状态
	}
}

