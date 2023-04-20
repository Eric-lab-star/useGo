package rss

type Deduper struct {
	s        Subscription
	updatesC chan Item
	closingC chan chan error
}

// Dedupe converts a Subscription that may send duplicate Items into
// one that doesn't.
func Dedupe(s Subscription) Subscription {
	d := &Deduper{
		s:        s,
		updatesC: make(chan Item),
		closingC: make(chan chan error),
	}
	go d.loop()
	return d
}

func (d *Deduper) loop() {
	in := d.s.Updates() // enable receive
	var pending Item
	var out chan Item // disable send
	seen := make(map[string]bool)
	for {
		select {
		case it := <-in:
			if !seen[it.GUID] {
				pending = it
				in = nil         // disable receive
				out = d.updatesC // enable send
				seen[it.GUID] = true
			}
		case out <- pending:
			in = d.s.Updates() // enable receive
			out = nil          // disable send
		case errc := <-d.closingC:
			err := d.s.Close()
			errc <- err
			close(d.updatesC)
			return
		}
	}
}

func (d *Deduper) Close() error {
	errc := make(chan error)
	d.closingC <- errc
	return <-errc
}

func (d *Deduper) Updates() <-chan Item {
	return d.updatesC
}
