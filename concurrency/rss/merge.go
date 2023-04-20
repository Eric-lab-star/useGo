package rss

type merge struct {
	subs    []Subscription
	updates chan Item
	quit    chan struct{}
	errs    chan error
}

// Merge returns a Subscription that merges the item streams from subs.
// Closing the merged subscription closes subs.
func Merge(subs ...Subscription) Subscription {

	m := &merge{
		subs:    subs,
		updates: make(chan Item),
		quit:    make(chan struct{}),
		errs:    make(chan error),
	}

	for _, sub := range subs {
		go func(s Subscription) {
			for {
				var it Item
				select {
				case it = <-s.Updates():
				case <-m.quit: // HL
					m.errs <- s.Close() // HL
					return              // HL
				}
				select {
				case m.updates <- it:
				case <-m.quit: // HL
					m.errs <- s.Close() // HL
					return              // HL
				}
			}
		}(sub)
	}

	return m
}

func (m *merge) Updates() <-chan Item {
	return m.updates
}

func (m *merge) Close() (err error) {
	close(m.quit) // HL
	for range m.subs {
		if e := <-m.errs; e != nil { // HL
			err = e
		}
	}
	close(m.updates) // HL
	return
}
