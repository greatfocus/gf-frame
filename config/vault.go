package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/greatfocus/gf-sframe/crypt"
)

// Vault struct
type Vault struct {
	URL         string `json:"url"`
	Application string `json:"application"`
	Impl        string `json:"impl"`
	Env         string `json:"env"`
}

// GetConfig method gets configf from vault
func (v *Vault) GetConfig(file string) Config {
	// read vault config
	var val = v.read(file)

	request := Vault{
		Application: val.Application,
		Impl:        val.Impl,
		Env:         os.Args[1],
	}
	reqBody, err := json.Marshal(request)
	if err != nil {
		log.Fatal(fmt.Println("Failed to get Vault config", err))
	}
	if err != nil {
		log.Fatal(fmt.Println("Failed to get Vault config", err))
	}

	client := http.Client{
		Timeout: time.Minute * 3,
		Transport: &http.Transport{
			TLSClientConfig: crypt.TLSClientConfig(),
		},
	}

	// make API call to vault
	resp, err := client.Post(val.URL, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		log.Fatal(fmt.Println("Failed to get Vault config"))
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(fmt.Println("Failed to get Vault config"))
	}

	// marshal te response
	var config Config
	err = json.Unmarshal(body, &config)
	if err != nil {
		log.Fatal(fmt.Println("Failed to get Vault config"))
	}

	// verify response
	if config.Impl == "" {
		log.Fatal(fmt.Println("Failed to get Vault config"))
	}

	// validate
	config.validate()

	return config
}

func (v *Vault) read(file string) Vault {
	log.Println("Reading configuration file")
	if len(file) < 1 {
		log.Fatal(fmt.Println("config file name is empty"))
	}

	jsonFile, err := os.OpenFile(file, os.O_RDONLY, 0600)
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Fatal(fmt.Println("cannot find location of config file", err))
	}

	// read the config file
	byteContent, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(fmt.Println("invalid config format", err))
	}

	// convert the config file bytes to json
	var result = Vault{}
	err = json.Unmarshal([]byte(byteContent), &result)
	if err != nil {
		log.Fatal(fmt.Println("Invalid config format", err))
	}

	// the closing of our jsonFile so that we can parse it later on
	_ = jsonFile.Close()

	return result
}
