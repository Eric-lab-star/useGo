package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/Eric-lab-star/useGo/concurrency/rss"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}
func main() {

	// Subscribe to some feeds, and create a merged update stream.
	merged := NaiveMerge(
		Subscribe(Fetch("blog.golang.org")),
		Subscribe(Fetch("googleblog.blogspot.com")),
		Subscribe(Fetch("googledevelopers.blogspot.com")),
	)

	// Close the subscriptions after some time.
	time.AfterFunc(3*time.Second, func() {
		fmt.Println("closed:", merged.Close())
	})

	// Print the stream.
	for it := range merged.Updates() {
		fmt.Println(it.Channel, it.Title)
	}

	fmt.Println("Number of Goroutines: ", runtime.NumGoroutine())
	panic("show goroutine stack")
}

// Fetch returns a Fetcher for Items from domain.
func Fetch(domain string) rss.Fetcher {
	return &rss.FakeFetcher{Channel: domain}
}

// Subscribe returns a new Subscription that uses fetcher to fetch Items.
func Subscribe(fetcher rss.Fetcher) rss.Subscription {
	s := &rss.Sub{
		Fetcher: fetcher,
	}
	go s.LoopCloseOnly()
	return s
}

func NaiveMerge(subs ...rss.Subscription) rss.Subscription {
	m := &rss.NaiveMerge{
		Subs: subs,
	}

	for _, sub := range subs {
		go func(s rss.Subscription) {
			for it := range s.Updates() {
				m.UpdatesC <- it // HL
			}
		}(sub)
	}

	return m
}
