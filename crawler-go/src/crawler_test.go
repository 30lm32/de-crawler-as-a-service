package main

import (
	"testing"
	"fmt"
	"log"
	"os"
	"bufio"
	"encoding/csv"
	"strconv"
)

func readLinks()(links map[link]int64){

	links = make(map[link]int64)

	f, err := os.Open("./test/unique_links.csv")

	if err != nil {
		log.Fatal(err.Error())
	}

	r := csv.NewReader(bufio.NewReader(f))
	records, err := r.ReadAll()

	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(records); i++ {
		record := records[i]
		n, err:= strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			log.Fatal(err.Error())
		}
		url0 := record[1]
		url1 := record[2]

		lnk := link{url0, url1}
		links[lnk] = n
	}

	return links
}

func readUrls()(urls map[string]int64){

	urls = make(map[string]int64)

	f, err := os.Open("./test/unique_urls.csv")

	if err != nil {
		log.Fatal(err.Error())
	}

	r := csv.NewReader(bufio.NewReader(f))
	records, err := r.ReadAll()

	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(records); i++ {
		record := records[i]
		n, err:= strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			log.Fatal(err.Error())
		}
		url := record[1]
		urls[url] = n
	}

	return urls
}

var urls map[string]int64
var links map[link]int64

func TestMain(m *testing.M){

	os.Setenv("TEST", "enable")

	urls = readUrls()
	links = readLinks()
	code := m.Run()

	os.Unsetenv("TEST")
	os.Exit(code)
}

func test(configFile string, t *testing.T) {

	var util Util

	config, err := util.ReadConfig(configFile)
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

	for k, ev := range urls {
		av := crawler.agent.getNumberOfUrls(k)
		if ev != av {
			t.Fatalf("Expected %d but got %d %s", ev, av, configFile)
		}
	}

	for k, ev := range links {

		av := crawler.agent.getNumberOfLinks(k)
		if ev != av {
			t.Fatalf("Expected %d but got %d %s", ev, av, configFile)
		}
	}

	ev := len(urls) + 1
	av := crawler.agent.getNumberOfUniqueUrls()
	if ev != av {
		t.Fatalf("Expected %d but got %d %s", ev, av, configFile)
	}

	//crawler.Reset()
	crawler.Close()
}

func TestCrawler_CrawlBFS(t *testing.T) {
	test("./test/config.test.bfs.json", t)
}

func TestCrawler_CrawlDFS(t *testing.T) {
	test("./test/config.test.dfs.json", t)
}