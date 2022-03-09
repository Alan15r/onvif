package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/use-go/onvif"
)

func main() {
	//создаем девайс(камеру) и подключаемся к ней
	device, err := onvif.NewDevice(onvif.DeviceParams{Xaddr: "172.22.226.8:10080"})
	if err != nil {
		panic(err)
	}

	//ptz
	//device
	//analytics
	//imaging
	//media

	//делаем запрос по нужному эндпоинту
	resp, err := device.CallMethod("media")
	if err != nil {
		log.Print("ERROR: ", err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("ERROR2: ", err.Error())
	}

	fmt.Println(string(body))
}
