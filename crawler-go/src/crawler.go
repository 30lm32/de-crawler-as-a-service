package main

import (
	"golang.org/x/net/html"
	"net/http"
	"log"
	"fmt"
	"time"
	"strings"
	"os"
)

type Crawler struct {
	agent  *Agent
	testOk bool //this flag added for testing purpose on local
}

func (crawler *Crawler) getHref(token html.Token, filter string) (bool, string) {
	for _, attr := range token.Attr {

		if attr.Key == keyHref{
			if  crawler.testOk {
				return true, attr.Val
			} else {
				if strings.Contains(attr.Val, filter) &&
					strings.Index(attr.Val, keyHttp) == 0 {

					return true, attr.Val
				}
			}
		}

	}
	return false, ""
}



func (crawler *Crawler) fetchUrl(pUrl string, chLink chan link, filterUrl string, method string) {
	resp, err := http.Get(pUrl)

	if err != nil {
		log.Println("Failed to crawl \"" + pUrl + "\"")
		return
	}

	body := resp.Body
	defer body.Close()

	tokenizer := html.NewTokenizer(body)

	isFirstUrlInBody := false
	for {
		tokenType := tokenizer.Next()

		switch {
		case tokenType == html.ErrorToken:
			return
		case tokenType == html.StartTagToken:

			token := tokenizer.Token()
			if token.Data != keyA {
				continue
			}

			ok, sUrl := crawler.getHref(token, filterUrl)
			if ok {

				if !crawler.agent.UrlExists(sUrl) {

					chLink <- link{pUrl, sUrl}
					switch method {
					case keyDFS:
						if !isFirstUrlInBody {
							isFirstUrlInBody = true
							crawler.fetchUrl(sUrl, chLink, filterUrl, method)
						} else {
							go crawler.fetchUrl(sUrl, chLink, filterUrl, method)
						}
						break
					case keyBFS:
						break
					}
				} else {
					chLink <- link{pUrl, sUrl}
				}
			}
		}
	}
}

func (crawler *Crawler) Crawl() {

	agent := crawler.agent
	config := agent.config
	method := config.Crawler_Method

	_, ok := os.LookupEnv("TEST")
	crawler.testOk  = ok

	agent.ResetByConfig()

	chLink := make(chan link)

	go crawler.fetchUrl(config.Target_Url, chLink, config.Filter_Url, method)
	agent.CountUrl(config.Target_Url)

	cont := true
	for cont {
		select {
		case lnk := <-chLink:

			if method == keyBFS {
				if !agent.UrlExists(lnk.SUrl) {
					go crawler.fetchUrl(lnk.SUrl, chLink, config.Filter_Url, method)
					log.Printf("%s\n", lnk.SUrl)
				}
			}

			agent.CountUrl(lnk.SUrl)
			agent.CountLink(lnk)

		case <-time.After(time.Second * 10):
			log.Println("Timeout: Since the crawler can not find no more new url during 10s, the crawler were timeout")
			cont = false
		}

	}

	log.Println(fmt.Sprintf("\nFound %d unique urls\n", agent.getNumberOfUniqueUrls()))
	log.Println(fmt.Sprintf("\nFound %d unique links\n", agent.getNumberOfUniqueLinks()))

	defer close(chLink)
}

func (crawler *Crawler) Close() {
	crawler.agent.Close()
}

func (crawler *Crawler) Reset() {
	crawler.agent.Reset()
}