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

// Impl struct
type Impl struct {
	Vault       string            `json:"vault"`
	Application string            `json:"application"`
	Impl        string            `json:"impl"`
	Env         string            `json:"env"`
	Scripts     map[string]string `json:"scripts"`
}

// GetConfig method gets configf from impl
func (i *Impl) GetConfig() Config {
	request := Impl{
		Application: i.Application,
		Impl:        i.Impl,
		Env:         i.Env,
	}
	reqBody, err := json.Marshal(request)
	if err != nil {
		log.Fatal(fmt.Println("Failed to get Impl config", err))
	}
	if err != nil {
		log.Fatal(fmt.Println("Failed to get Impl config", err))
	}

	client := http.Client{
		Timeout: time.Minute * 3,
		Transport: &http.Transport{
			TLSClientConfig: crypt.TLSClientConfig(),
		},
	}

	// make API call to impl
	resp, err := client.Post(i.Vault, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		log.Fatal(fmt.Println("Failed to get Impl config"))
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(fmt.Println("Failed to get Impl config"))
	}

	// marshal te response
	var config Config
	err = json.Unmarshal(body, &config)
	if err != nil {
		log.Fatal(fmt.Println("Failed to get Impl config"))
	}

	// verify response
	if config.Impl == "" {
		log.Fatal(fmt.Println("Failed to get Impl config"))
	}

	// validate
	config.validate()

	return config
}

func (i *Impl) read(file string) Impl {
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
	var result = Impl{}
	err = json.Unmarshal([]byte(byteContent), &result)
	if err != nil {
		log.Fatal(fmt.Println("Invalid config format", err))
	}

	// the closing of our jsonFile so that we can parse it later on
	_ = jsonFile.Close()

	return result
}
