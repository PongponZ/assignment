package model

type Source struct {
	Username  string `json:"username"`
	IpAddress string `json:"ip_address"`
}

type Destination struct {
	Username  string  `json:"username"`
	IpAddress string  `json:"ip_address"`
	Address   Address `json:"address"`
}

type Address struct {
	CountryCode string `json:"country_code"`
	CountryName string `json:"country_name"`
}
