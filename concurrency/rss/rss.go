package rss

import "time"

// An Item is a stripped-down RSS item.
type Item struct{ Title, Channel, GUID string }

type Fetcher interface {
	Fetch() (items []Item, next time.Time, err error)
}

// A Subscription delivers Items over a channel.  Close cancels the
// subscription, closes the Updates channel, and returns the last fetch error,
// if any.
type Subscription interface {
	Updates() <-chan Item
	Close() error
}
