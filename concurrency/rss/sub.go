package rss

import "time"

// Sub implements the Subscription interface.
type Sub struct {
	Fetcher  Fetcher         // fetches items
	UpdatesC chan Item       // sends items to the user
	ClosingC chan chan error // for Close
}

func (s *Sub) Updates() <-chan Item {
	return s.UpdatesC
}

func (s *Sub) Close() error {
	errc := make(chan error)
	s.ClosingC <- errc // HLchan
	return <-errc      // HLchan
}

// loopCloseOnly is a version of loop that includes only the logic
// that handles Close.
func (s *Sub) LoopCloseOnly() {
	var err error // set when Fetch fails
	for {
		select {
		case errc := <-s.ClosingC: // HLchan
			errc <- err       // HLchan
			close(s.UpdatesC) // tells receiver we're done
			return
		}
	}

}

// loopFetchOnly is a version of loop that includes only the logic
// that calls Fetch.
func (s *Sub) LoopFetchOnly() {
	var (
		pending []Item    // appended by fetch; consumed by send
		next    time.Time // initially January 1, year 0
		err     error
	)
	for {
		var fetchDelay time.Duration // initally 0 (no delay)
		if now := time.Now(); next.After(now) {
			fetchDelay = next.Sub(now)
		}
		startFetch := time.After(fetchDelay)

		select {
		case <-startFetch:
			var fetched []Item
			fetched, next, err = s.Fetcher.Fetch()
			if err != nil {
				next = time.Now().Add(10 * time.Second)
				break
			}
			pending = append(pending, fetched...)
		}
	}

}

// loopSendOnly is a version of loop that includes only the logic for
// sending items to s.updates.
func (s *Sub) LoopSendOnly() {

	var pending []Item // appended by fetch; consumed by send
	for {
		var first Item
		var updates chan Item // HLupdates
		if len(pending) > 0 {
			first = pending[0]
			updates = s.UpdatesC // enable send case // HLupdates
		}

		select {
		case updates <- first:
			pending = pending[1:]
		}
	}

}

// mergedLoop is a version of loop that combines loopCloseOnly,
// loopFetchOnly, and loopSendOnly.
func (s *Sub) MergedLoop() {

	var (
		pending []Item
		next    time.Time
		err     error
	)

	for {
		var fetchDelay time.Duration
		if now := time.Now(); next.After(now) {
			fetchDelay = next.Sub(now)
		}
		startFetch := time.After(fetchDelay)

		var first Item
		var updates chan Item
		if len(pending) > 0 {
			first = pending[0]
			updates = s.UpdatesC // enable send case
		}

		select {
		case errc := <-s.ClosingC: // HLcases
			errc <- err
			close(s.UpdatesC)
			return

		case <-startFetch: // HLcases
			var fetched []Item
			fetched, next, err = s.Fetcher.Fetch() // HLfetch
			if err != nil {
				next = time.Now().Add(10 * time.Second)
				break
			}
			pending = append(pending, fetched...) // HLfetch

		case updates <- first: // HLcases
			pending = pending[1:]
		}

	}
}
