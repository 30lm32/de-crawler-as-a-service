package main

import (
	"testing"
	"fmt"
)

var testConfig = &configuration{
	App_Name: "test-crawler",
	Redis_Address: "redis:6379",
	Neo4j_Bolt_Address: "bolt://neo4j:7687",
	Redis_Connection_Attempts: 15,
	Neo4j_Connection_Attempts: 15,
	Attempts_Time_Interval: "1000ms",
	Reset_Redis: true,
	Reset_Neo4j: true,
}

//func TestMain(m *testing.M) {
//	setup()
//	code := m.Run()
//	shutdown()
//	os.Exit(code)
//}

func TestNewAgentNilConfig(t *testing.T) {

	var config configuration
	agent := NewAgent(&config)

	if agent != nil {
		t.Fatalf("Expected %s but got %s", nil, agent)
	}
}

func TestNewAgentNotNilConfig(t *testing.T) {

	agent := NewAgent(testConfig)
	if agent == nil {
		t.Fatalf("Expected %s but got %s", agent, nil)
	}

	pong, err := agent.client.Ping().Result()
	if err != nil {
		t.Fatalf("Expected %s but got %s", nil, err.Error())
	}

	if pong != "PONG" {
		t.Fatalf("Expected %s but got %s", "PONG", pong)
	}

	agent.Close()
}

func TestNewAgentArbitraryConfig(t *testing.T) {

	config := &configuration{
		Redis_Address: "",
		Neo4j_Bolt_Address: "",
		Redis_Connection_Attempts: 15,
		Neo4j_Connection_Attempts: 15,
		Attempts_Time_Interval: "",
	}

	agent := NewAgent(config)
	if agent != nil {
		t.Fatalf("Expected %s but got %s", nil, agent)
	}
}

func TestAgent_CountLink(t *testing.T) {

	agent := NewAgent(testConfig)
	if agent == nil {
		t.Fatalf("Expected %s but got %s", agent, nil)
	}
	agent.ResetByConfig()

	formatUrl := "test123://%s.%d.com"
	for i := 0 ; i < 12; i++ {
		m := i % 4
		pUrl := fmt.Sprintf(formatUrl, "p", m)
		sUrl := fmt.Sprintf(formatUrl, "s", m)
		lnk := link{pUrl,sUrl}
		agent.CountLink(lnk)
	}

	for i := 0; i < 4; i++ {
		pUrl := fmt.Sprintf(formatUrl, "p", i)
		sUrl := fmt.Sprintf(formatUrl, "s", i)
		lnk := link{pUrl,sUrl}
		actualCount := agent.getNumberOfLinks(lnk)
		if actualCount != 3 {
			t.Fatalf("Expected number of links %d but got %d", 3, actualCount)
		}
	}

	agent.Close()
}

func TestAgent_CountUrl(t *testing.T) {

	agent := NewAgent(testConfig)
	if agent == nil {
		t.Fatalf("Expected %s but got %s", agent, nil)
	}

	agent.ResetByConfig()

	url := "test123://www.%d.com"
	for i := 0 ; i < 12; i++ {
		m := i % 4
		url := fmt.Sprintf(url, m)
		agent.CountUrl(url)
	}

	for i := 0; i < 4; i++ {
		url := fmt.Sprintf(url, i)
		actualCount := agent.getNumberOfUrls(url)
		if actualCount != 3 {
			t.Fatalf("Expected number of links %d but got %d", 3, actualCount)
		}
	}

	agent.Close()
}

func TestAgent_UrlExists(t *testing.T) {

	agent := NewAgent(testConfig)
	if agent == nil {
		t.Fatalf("Expected %s but got %s", agent, nil)
	}

	agent.ResetByConfig()

	url := "test123://www.%d.com"
	for i := 0 ; i < 5; i++ {
		url := fmt.Sprintf(url, i)
		agent.CountUrl(url)
	}

	for i := 0; i < 5; i++ {
		url := fmt.Sprintf(url, i)
		actualResult := agent.UrlExists(url)
		if !actualResult {
			t.Fatalf("Expected result %b but got %d", true, actualResult)
		}
	}

	agent.Close()
}
