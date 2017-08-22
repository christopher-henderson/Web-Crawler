package extractor

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func ExtractAll(URL string) ([]byte, []string) {
	resp, err := http.Get(URL)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	defer resp.Body.Close()
	U, _ := url.Parse(URL)
	domain := U.Hostname()
	scheme := U.Scheme
	var w []byte
	writer := bytes.NewBuffer(w)
	body := io.TeeReader(bufio.NewReader(resp.Body), writer)
	tokenizer := html.NewTokenizer(body)
	var urls []string
	for {
		if next := tokenizer.Next(); next == html.ErrorToken {
			return writer.Bytes(), urls
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
