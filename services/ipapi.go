package services

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"assignment/model"
)

type IpApiResponse struct {
	CountryName string `json:"country_name"`
	CountryCode string `json:"country_code"`
}

func GetIpInfo(ip string) model.Address {
	url := "https://ipapi.co/" + ip + "/json/"

	client := http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "ipapi.co/#go-v1.5")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var info model.Address
	err = json.Unmarshal(body, &info)
	if err != nil {
		log.Fatal(err)
	}

	return info
}
