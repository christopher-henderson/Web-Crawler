package frontier

import (
	"sync"
	"time"
)

var maxFrontier = 10000
var queue = make(chan string, maxFrontier)
var seen = make(map[string]bool)
var lock = sync.Mutex{}
var seed = []string{"http://pornhub.com"}

func init() {
	Push(seed)
}

func Push(urls []string) {
	lock.Lock()
	defer lock.Unlock()
	for _, url := range urls {
		if _, ok := seen[url]; ok {
			// log.Println("Duplicate URL")
			continue
		}
		seen[url] = true
		// log.Printf("Pushing : %v\n", url)
		select {
		case queue <- url:
		default:
		}
	}
}

func Pop() (string, bool) {
	select {
	case <-time.After(time.Second * 5):
		return "", false
	case url := <-queue:
		return url, true
	}
}
