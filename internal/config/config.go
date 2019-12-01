package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// Client handler
type Client struct {
	Mac                 string `json:"mac"`
	LocalProxyBind      string `json:"local_proxy_bind"`
	RemoteProxyData     string `json:"remote_proxy_data"`
	RemoteProxyControl  string `json:"remote_proxy_control"`
	LocalBindMin        uint16 `json:"local_bind_min"`
	LocalBindMax        uint16 `json:"local_bind_max"`
	Browser             string `json:"browser"`
	ProxyScheme         string `json:"proxy_scheme"`
	ProxyTLSCertificate string `json:"proxy_tls_certificate"`
	ProxyTLSKey         string `json:"proxy_tls_key"`
}

// Server handler
type Server struct {
	ControlBind string `json:"control_bind"`
	DataBind    string `json:"data_bind"`
	DbDriver    string `json:"db_driver"`
	DbConn      string `json:"db_conn"`
}

// ReadConfigClient reads the client configuration
func ReadConfigClient(filename string) Client {
	var config Client

	json.Unmarshal(readConfig(filename), &config)

	return config
}

// ReadConfigServer reads the server configuration
func ReadConfigServer(filename string) Server {
	var config Server

	json.Unmarshal(readConfig(filename), &config)

	return config
}

func readConfig(filename string) []byte {
	jsonFile, err := os.Open(filename)

	if err != nil {
		log.Fatalln("Cannot open file:", err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatalln("Cannot read file:", err)
	}

	return byteValue
}
