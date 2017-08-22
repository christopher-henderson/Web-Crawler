package main

import (
	"crawler/extractor"
	"fmt"
	"net/http"
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
		fmt.Println(err)
	} else {
		fmt.Println(resp)
	}
	ch <- true
}

func main() {
	seed := "http://reuters.com"
	fmt.Println("DO IT")
	start := time.Now()
	extractor.ExtractAll(seed)
	fmt.Printf("%v\n", time.Now().Sub(start))
	// fmt.Println(string(content))
	// fmt.Println(urls)
}
