package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func test(url string) bool {
	resp, err := http.Head(url)
	if err != nil {
		fmt.Println(err)
		return false
	}
	header := resp.Header.Get("Content-Type")
	if strings.HasPrefix(header, "text/html") {
		return true
	}
	return false
}

func crawl(URL string) []string {
	resp, err := http.Get(URL)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	U, _ := url.Parse(URL)
	domain := U.Hostname()
	scheme := U.Scheme
	tokenizer := html.NewTokenizer(resp.Body)
	var urls []string
	for {
		if next := tokenizer.Next(); next == html.ErrorToken {
			return urls
		} else if next != html.StartTagToken {
			continue
		} else if token := tokenizer.Token(); token.Data != "a" {
			continue
		} else {
			if target, err := extractURL(domain, scheme, token.Attr[0].Val); err == nil {
				urls = append(urls, target)
			} else {
				fmt.Println(err)
			}
		}
	}
}

func extractURL(domain string, scheme, href string) (string, error) {
	URL, err := url.Parse(href)
	if err != nil {
		return "", err
	}
	if URL.IsAbs() {
		return URL.String(), nil
	}
	URL = &url.URL{}
	URL.Host = domain
	URL.Scheme = scheme
	URL.Path = string(href)
	return URL.String(), nil
}

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
		fmt.Println(err)
	} else {
		fmt.Println(resp)
	}
	ch <- true
}

func main() {
	seed := "http://reuters.com"
	if !test(seed) {
		fmt.Println("LOL NO")
		return
	}
	fmt.Println("DO IT")
	urls := crawl(seed)
	validate(urls)
	// fmt.Println(urls)
}
