package main

import (
	"fmt"
	"os"
	"encoding/json"
)

type configuration struct {
	App_Name                  string
	Redis_Address             string
	Neo4j_Bolt_Address        string
	Target_Url                string
	Filter_Url                string
	Reset_Redis               bool
	Reset_Neo4j               bool
	Redis_Connection_Attempts int64
	Neo4j_Connection_Attempts int64
	Attempts_Time_Interval    string
	Crawler_Method            string
}

const (
	keyUrl = "%s_url"
	keyLink = "%s_link"
	formatLink = "%s->%s"

	keyHref = "href"
	keyA = "a"
	keyHttp = "http"

	keyDFS = "dfs"
	keyBFS = "bfs"
)

type link struct {
	PUrl string
	SUrl string
}

type message struct {
	holder interface{}
	err    error
}

type Util struct {}

func (util *Util) validateConfig(config *configuration) error {
	if config.Redis_Address == "" {
		return fmt.Errorf("Please, assign Redis server address!\n")
	}

	if config.Target_Url == "" {
		return fmt.Errorf("Please, enter 'target_url' in config file\n")
	}

	if config.Filter_Url == "" {
		return fmt.Errorf("Please, enter 'filter_url' in config file\n")
	}

	if config.Attempts_Time_Interval == "" {
		return fmt.Errorf("Please, enter 'time_interval_unit' in config file\n")
	}

	if config.Crawler_Method == "" {
		return fmt.Errorf("Please, enter 'crawler_method' in config file\n")
	}

	return nil
}

func (util *Util) ReadConfig(filename string) (*configuration, error) {

	var config configuration

	file, _ := os.Open(filename)
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&config)

	if err != nil {
		return nil, err
	}

	if &config == nil {
		return nil, fmt.Errorf("Configuration is not created!")
	}

	err = util.validateConfig(&config)
	if err != nil {
		return nil, err
	} else {
		return &config, nil
	}
}





