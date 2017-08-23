package main

import (
	"crawler/extractor"
	"log"
	"net/http"
	"net/url"
	"time"
)

func validate(urls []string) {
	ch := make(chan bool)
	count := 0
	for _, u := range urls {
		go v(u, ch)
		count++
	}
	for ; count > 0; count-- {
		<-ch
	}
}

func v(u string, ch chan bool) {
	if resp, err := http.Head(u); err != nil {
		log.Println(err)
	} else {
		if resp.StatusCode != 200 {
			log.Printf("%v: %v\n", resp.StatusCode, u)
		}
	}
	ch <- true
}

func countDomains(urls []string) int {
	domains := make(map[string]bool)
	for _, u := range urls {
		U, _ := url.Parse(u)
		domains[U.Host] = true
	}
	log.Println(domains)
	return len(domains)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	seed := "http://foxnews.com"
	log.Println("DO IT")
	start := time.Now()
	_, urls := extractor.ExtractAll(seed)
	// log.Println(string(content))
	log.Println(urls)
	validate(urls)
	log.Printf("%v\n", time.Now().Sub(start))
	log.Println(len(urls))
	log.Println(countDomains(urls))
}
