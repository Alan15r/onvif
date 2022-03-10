package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/use-go/onvif"
	"github.com/use-go/onvif/media"
	"github.com/yakovlevdmv/goonvif/Media"
	onv "github.com/yakovlevdmv/goonvif/xsd/onvif"
)

type Envelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    Body     `xml:"Body"`
}

type Body struct {
	GetProfilesResponse media.GetProfilesResponse `xml:"GetProfilesResponse"`
}

type EnvelopeShot struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    BodyShot `xml:"Body"`
}

type BodyShot struct {
	GetSnapshotUriResponse media.GetSnapshotUriResponse `xml:"GetSnapshotUriResponse"`
}

func main() {
	//создаем девайс(камеру) и подключаемся к ней
	xaddr := "172.22.226.8:10080"
	device, err := onvif.NewDevice(onvif.DeviceParams{Xaddr: xaddr})
	if err != nil {
		panic(err)
	}

	//получаем профили
	getProf, err := GetProfiles(device)
	if err != nil {
		log.Fatal("error get profiles: ", err.Error())
	}
	profiles := getProf.Profiles

	profToken := string(profiles[0].Name)
	refToken := onv.ReferenceToken(profToken)

	//получаем uri для снимка
	uri, err := GetUri(device, refToken)
	if err != nil {
		log.Fatal("error get uri: ", err.Error())
	}

	//переделываем адрес
	url := remake(uri, xaddr)

	//запрашиваем кадр
	shot, err := GetShot(url)
	if err != nil {
		log.Fatal("error get shot: ", err.Error())
	}

	fmt.Println(shot)
}

func GetShot(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func GetUri(device *onvif.Device, refToken onv.ReferenceToken) (string, error) {
	getUri := Media.GetSnapshotUri{ProfileToken: refToken}
	resp, err := device.CallMethod(getUri)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var res EnvelopeShot
	err = xml.Unmarshal(body, &res)
	if err != nil {
		return "", err
	}

	uri := string(res.Body.GetSnapshotUriResponse.MediaUri.Uri)
	return uri, nil
}

func GetProfiles(device *onvif.Device) (media.GetProfilesResponse, error) {
	getProf := Media.GetProfiles{}
	resp, err := device.CallMethod(getProf)
	if err != nil {
		return media.GetProfilesResponse{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return media.GetProfilesResponse{}, err
	}

	var res Envelope
	err = xml.Unmarshal(body, &res)
	if err != nil {
		return media.GetProfilesResponse{}, err
	}

	return res.Body.GetProfilesResponse, nil
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
