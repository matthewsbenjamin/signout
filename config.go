package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// Configs stores the database credentials
type Configs struct {
	Port string `yaml:"port"`
}

type dbCred struct {
	User     string `yaml:"user"`
	Pwd      string `yaml:"pwd"`
	DB       string `yaml:"db"`
	Endpoint string `yaml:"endpoint"`
	Port     string `yaml:"port"`
	DBname   string `yaml:"dbname"`
}

// Creds contains various credential information - primarily database
type Creds struct {
	DB dbCred `yaml:"db"`
}

func (c *Configs) getConf() *Configs {

	file, err := ioutil.ReadFile("config/config.yaml")
	if err != nil {
		log.Println(err)
	}
	err = yaml.Unmarshal(file, c)
	if err != nil {
		log.Println(err)
	}

	return c
}

func (c *Creds) getCred() *Creds {

	file, err := ioutil.ReadFile("config/cred.yaml")
	if err != nil {
		log.Println(err)
	}
	err = yaml.Unmarshal(file, c)
	if err != nil {
		log.Println(err)
	}

	return c
}

func (c *Creds) dbCred() string {

	c.getCred()

	s := fmt.Sprintf("%s:%s@tcp(%s)/%s", c.DB.User, c.DB.Pwd, c.DB.Endpoint, c.DB.DBname)

	return s
}
