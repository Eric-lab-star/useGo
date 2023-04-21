package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"
)

func main() {
	m := Merge(
		Subscribe(Fetch("googl.com")),
		Subscribe(Fetch("blog.com")),
		Subscribe(Fetch("news.com")),
	)
	time.AfterFunc(1*time.Second, func() {
		fmt.Println("closing", m.Close())
	})

	for i := range m.Updates() {
		fmt.Println(i.Title)
		fmt.Println(i.Channel, i.GUID)
	}

	pprof.Lookup("goroutine").WriteTo(os.Stderr, 1)

}

type Item struct{ Title, Channel, GUID string }

type Fetcher interface {
	fetch() (items []Item, next time.Time, err error)
}

type FakeFetcher struct {
	Channel string
	Items   []Item
}

var FakeDuplicates bool

func (f *FakeFetcher) fetch() (items []Item, next time.Time, err error) {
	now := time.Now()
	next = now.Add(time.Duration(rand.Intn(5)) * 500 * time.Millisecond)
	item := Item{
		Channel: f.Channel,
		Title:   fmt.Sprintf("Items %d", len(f.Items)),
	}
	item.GUID = item.Channel + "/" + item.Title
	f.Items = append(f.Items, item)
	if FakeDuplicates {
		items = f.Items
	} else {
		items = []Item{item}
	}
	return
}

func Fetch(domain string) Fetcher {
	return &FakeFetcher{Channel: domain}
}

type sub struct {
	fetcher Fetcher
	updates chan Item
	closing chan chan error
}

type subscription interface {
	Updates() <-chan Item
	Close() error
}

func (s *sub) Updates() <-chan Item {
	return s.updates
}

func (s *sub) Close() error {
	errc := make(chan error)
	s.closing <- errc
	return <-errc
}

func (s *sub) MergedLoop() {

	var (
		pending []Item
		next    time.Time
		err     error
	)

	for {

		var (
			fetchDelay time.Duration
			first      Item
			updates    chan Item
		)
		if now := time.Now(); next.After(now) {
			fetchDelay = next.Sub(now)
		}

		if len(pending) > 0 {
			first = pending[0]
			updates = s.updates
		}

		startFetch := time.After(fetchDelay)

		select {
		case errc := <-s.closing:
			errc <- err
			close(s.updates)
			return
		case <-startFetch:
			var fetched []Item
			fetched, next, err = s.fetcher.fetch()
			if err != nil {
				next = time.Now().Add(10 * time.Second)
				break
			}
			pending = append(pending, fetched...)
		case updates <- first:
			pending = pending[1:]
		}

	}
}

type merged struct {
	subs    []subscription
	updates chan Item
}

func (m *merged) Updates() <-chan Item {
	return m.updates
}

func (m *merged) Close() (err error) {
	for _, sub := range m.subs {
		if e := sub.Close(); err == nil && e != nil {
			err = e
		}

	}
	close(m.updates)
	return
}

func Subscribe(fetcher Fetcher) subscription {
	s := &sub{
		fetcher: fetcher,
		updates: make(chan Item),
		closing: make(chan chan error),
	}
	go s.MergedLoop()
	return s
}

func Merge(subs ...subscription) subscription {
	m := &merged{
		subs:    subs,
		updates: make(chan Item),
	}

	for _, sub := range subs {
		go func(sub subscription) {
			for i := range sub.Updates() {
				m.updates <- i
			}
		}(sub)
	}

	return m
}
