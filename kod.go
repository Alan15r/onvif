package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/use-go/onvif"
	"github.com/yakovlevdmv/goonvif/Media"
	onv "github.com/yakovlevdmv/goonvif/xsd/onvif"
)

type Envelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    Body     `xml:"Body"`
}

type Body struct {
	GetProfilesResponse GetProfilesResponse `xml:"GetProfilesResponse"`
}

type GetProfilesResponse struct {
	Profiles []Profile `xml:"Profiles"`
}

type Profile struct {
	Name string `xml:"Name"`
}

//

type EnvelopeShot struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    BodyShot `xml:"Body"`
}

type BodyShot struct {
	GetSnapshotUriResponse GetSnapshotUriResponse `xml:"GetSnapshotUriResponse"`
}

type GetSnapshotUriResponse struct {
	MediaUri MediaUri `xml:"MediaUri"`
}

type MediaUri struct {
	Uri string `xml:"Uri"`
}

func main() {
	//создаем девайс(камеру) и подключаемся к ней
	xaddr := "172.22.226.8:10080"
	device, err := onvif.NewDevice(onvif.DeviceParams{Xaddr: xaddr})
	if err != nil {
		panic(err)
	}

	//получаем профили
	getProf := Media.GetProfiles{}
	resp, err := device.CallMethod(getProf)
	if err != nil {
		log.Println("ERROR.0", err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("ERROR2: ", err.Error())
	}

	var res Envelope
	err = xml.Unmarshal(body, &res)
	if err != nil {
		panic(err)
	}

	profToken := res.Body.GetProfilesResponse.Profiles[0].Name
	refToken := onv.ReferenceToken(profToken)

	//получаем uri для снимка
	getUri := Media.GetSnapshotUri{ProfileToken: refToken}
	resp2, err := device.CallMethod(getUri)
	if err != nil {
		log.Fatal("ERROR resp2: ", err.Error())
	}

	body2, err := ioutil.ReadAll(resp2.Body)
	if err != nil {
		log.Fatal("ERROR body2: ", err.Error())
	}

	var res1 EnvelopeShot
	err = xml.Unmarshal(body2, &res1)
	if err != nil {
		panic(err)
	}

	uri := res1.Body.GetSnapshotUriResponse.MediaUri.Uri

	//переделываем адрес
	url := remake(uri, xaddr)
	fmt.Println(url)

	//запрашиваем кадр
	resp3, err := http.Get(url)
	if err != nil {
		log.Fatal("ERROR resp3: ", err.Error())
	}

	body3, err := ioutil.ReadAll(resp3.Body)
	if err != nil {
		log.Fatal("ERROR body3: ", err.Error())
	}

	fmt.Println(body3)
}

func remake(uri, xaddr string) string {
	url := uri[:7] + xaddr
	for i, v := range uri {
		if v == '/' && i > 7 {
			url += uri[i:]
			break
		}
	}

	return url
}
