package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

//
type Config struct {
	Port string
}

type dbCred struct {
	User     string
	Pwd      string
	Endpoint string
	Port     string
	DBname   string
}

// Creds contains various credential information - primarily database
type Creds struct {
	DB dbCred
}

// GetConfig this will return a struct of type config
func GetConfig() Config {

	var i interface{}

	c := getYaml(i, "config/config.yaml")

	C, _ := c.(Config)

	return C

}

// GetCreds will return a struct of creds
func GetCreds() Creds {

	var i interface{}

	c := getYaml(i, "config/creds.yaml")

	C, _ := c.(Creds)

	return C
}

func getYaml(s interface{}, f string) interface{} {

	// s struct
	// f file

	file, err := ioutil.ReadFile(f)
	if err != nil {
		log.Fatal("Configuration file failed to open correctly")
	}

	// unmarshal the file into the config struct
	err = yaml.Unmarshal([]byte(file), &s)
	if err != nil {
		log.Fatal("Configuration file failed to parse correctly")
	}

	return s

}
