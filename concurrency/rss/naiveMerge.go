package rss

type NaiveMerge struct {
	Subs     []Subscription
	UpdatesC chan Item
}

func (m *NaiveMerge) Close() (err error) {
	for _, sub := range m.Subs {
		if e := sub.Close(); err == nil && e != nil {
			err = e
		}
	}
	close(m.UpdatesC) // HL
	return
}

func (m *NaiveMerge) Updates() <-chan Item {
	return m.UpdatesC
}
