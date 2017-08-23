package extractor

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

type Map map[string]bool

func (m Map) Keys() []string {
	var keys []string
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

func ExtractAll(URL string) ([]byte, []string) {
	resp, err := http.Get(URL)
	if err != nil {
		log.Println(err)
		return nil, nil
	}
	defer resp.Body.Close()
	U, _ := url.Parse(URL)
	domain := U.Hostname()
	scheme := U.Scheme
	contentWriter := bytes.NewBuffer(make([]byte, 0))
	t := io.TeeReader(bufio.NewReader(resp.Body), contentWriter)
	tokenizer := html.NewTokenizer(t)
	urls := make(Map)
	for {
		if next := tokenizer.Next(); next == html.ErrorToken {
			return contentWriter.Bytes(), urls.Keys()
		} else if next != html.StartTagToken {
			continue
		} else if token := tokenizer.Token(); token.Data != "a" {
			continue
		} else {
			if target, err := extractURL(domain, scheme, token.Attr[0].Val); err == nil {
				urls[target] = true
			} else {
				log.Println(err)
			}
		}
	}
}

func extractURL(domain, scheme, href string) (string, error) {
	href = strings.TrimSpace(href)
	if unescapedHref, err := url.PathUnescape(href); err != nil {
		log.Println(err)
		return "", err
	} else {
		href = unescapedHref
	}
	if strings.HasPrefix(href, "//") {
		href = strings.Join([]string{scheme, "://", href[2:]}, "")
	}
	URL, err := url.Parse(href)
	if err != nil {
		return "", err
	}
	if URL.IsAbs() {
		u, _ := url.PathUnescape(URL.String())
		return u, nil
	}
	URL = &url.URL{}
	URL.Host = domain
	URL.Scheme = scheme
	URL.Path = href
	if result, err := url.PathUnescape(URL.String()); err != nil {
		return "", err
	} else {
		return result, nil
	}
}

func test(url string) bool {
	resp, err := http.Head(url)
	if err != nil {
		log.Println(err)
		return false
	}
	header := resp.Header.Get("Content-Type")
	if strings.HasPrefix(header, "text/html") {
		return true
	}
	return false
}
