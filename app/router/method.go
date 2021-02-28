package router

import (
	"encoding/base64"
	"errors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"oims/app/docker"
	"oims/error"
	"os"
	"strconv"
	"time"
)

type image struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func getJpeg(c *gin.Context) interface{} {
	images := new([]*image)
	err := c.BindJSON(images)
	if err != nil {
		panic(error.NewHttpError(403, "40301", err.Error()))
	}
	TimeStamp := time.Now().Unix()
	basePath := "../historys/"
	_, err = os.Stat(basePath)
	if err != nil {
		panic(error.NewHttpError(500, "50001", err.Error()))
	}
	filePath := basePath + strconv.FormatInt(TimeStamp, 10)
	err = os.Mkdir(filePath, os.ModePerm)
	if err != nil {
		panic(error.NewHttpError(500, "50002", err.Error()))
	}

	if len(*images) == 0 {
		panic(error.NewHttpError(403, "40301", "no images"))
	}

	for _, file := range *images {
		filename := filePath + "/" + file.Name
		f, err := os.Create(filename)
		if err != nil {
			panic(err)
		}
		defer func() {
			_ = f.Close()
		}()
		maxLen := base64.StdEncoding.DecodedLen(len([]byte(file.Data)))
		dst := make([]byte, maxLen)
		// n是实际大小
		n, err := base64.StdEncoding.Decode(dst, []byte(file.Data))
		if err != nil {
			panic(error.NewHttpError(403, "40302", err.Error()))
		}
		n, err = f.Write(dst[:n]) //新的n是实际写入
		if err != nil {
			panic(err)
		} else if n == 0 {
			panic(errors.New(file.Name + "write null"))
		}
	}

	// docker 处理图片
	docker.Run(strconv.FormatInt(TimeStamp, 10))

	return TimeStamp
}

func getXml(c *gin.Context) interface{} {
	timeStamp := c.Query("id")
	resultPath := "../results/" + timeStamp + ".xml"
	f, err := os.Open(resultPath)
	if os.IsNotExist(err) {
		panic(error.NewHttpError(404, "40401", resultPath+" is not exist"))
	}
	data, err := ioutil.ReadFile(f.Name())
	if err != nil {
		panic(err)
	}
	result := base64.StdEncoding.EncodeToString(data)
	return result
}
