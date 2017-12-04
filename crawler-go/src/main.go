package main

import (
	"os"
	"log"
	"fmt"
)

func main() {
	confFile := os.Args[1]

	var util Util

	config, err := util.ReadConfig(confFile)
	if err != nil{
		log.Fatalf(err.Error())
		return
	}

	agent := NewAgent(config)
	if agent == nil{
		log.Fatalf(fmt.Errorf("Agent was not created...").Error())
		return
	}

	crawler := Crawler{agent:agent}
	crawler.Crawl()

	crawler.Close()
}
