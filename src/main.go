package main

import (
	"crawler/extractor"
	"crawler/frontier"
	"crawler/repository"
	"log"
	"net/http"
	"net/url"
	"sync"
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

func work(in chan string, wg *sync.WaitGroup) {
	for url := range in {
		content, urls, ok := extractor.ExtractAll(url)
		if !ok {
			return
		}
		frontier.Push(urls)
		repository.Save(url, content)
	}
	wg.Done()
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	maxRecords := 1000
	maxWorkers := 20
	wg := &sync.WaitGroup{}
	urls := make(chan string)
	defer wg.Wait()
	defer close(urls)
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go work(urls, wg)
	}
	start := time.Now()
	for i := 0; i < maxRecords; i++ {
		if url, ok := frontier.Pop(); ok {
			urls <- url
		} else {
			break
		}
	}
	log.Println(time.Now().Sub(start))
}

// func main() {
// log.SetFlags(log.LstdFlags | log.Lshortfile)
// 	log.Println(frontier.Pop())
// 	frontier.IncrementVisited()
// 	go func() {
// 		// time.Sleep(time.Second * 3)
// 		frontier.Push([]string{"http://google.com"})
// 	}()
// 	log.Println(frontier.Pop())
// 	go func() {
// 		// time.Sleep(time.Second * 3)
// 		for i := 0; i < 1000; i++ {
// 			frontier.IncrementVisited()
// 		}
// 	}()
// 	log.Println(frontier.Pop())
// 	frontier.Close()
// }

// func main() {
// 	seed := "http://foxnews.com"
// 	log.Println("DO IT")
// 	start := time.Now()
// 	_, urls := extractor.ExtractAll(seed)
// 	// log.Println(string(content))
// 	log.Println(urls)
// 	validate(urls)
// 	log.Printf("%v\n", time.Now().Sub(start))
// 	log.Println(len(urls))
// 	log.Println(countDomains(urls))
// }
