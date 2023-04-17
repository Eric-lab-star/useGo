package rss

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
