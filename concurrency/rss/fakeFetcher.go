package rss

import (
	"fmt"
	"math/rand"
	"time"
)

type FakeFetcher struct {
	Channel string
	Items   []Item
}

// FakeDuplicates causes the fake fetcher to return duplicate items.
var FakeDuplicates bool

// fakeFetcher implements Fetcher interface
func (f *FakeFetcher) Fetch() (items []Item, next time.Time, err error) {
	now := time.Now()
	next = now.Add(time.Duration(rand.Intn(5)) * 500 * time.Millisecond)
	item := Item{
		Channel: f.Channel,
		Title:   fmt.Sprintf("Item %d", len(f.Items)),
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
