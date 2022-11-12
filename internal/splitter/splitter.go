package splitter

import "sync"

type SseSplit struct {
	m sync.Mutex
	s map[chan<- string]struct{}
}

func NewSseSplitter() *SseSplit {
	return &SseSplit{s: make(map[chan<- string]struct{})}
}

func (f *SseSplit) Subscribe(c chan<- string) {
	f.m.Lock()
	defer f.m.Unlock()

	f.s[c] = struct{}{}
}

func (f *SseSplit) Unsubsctibe(c chan<- string) {
	f.m.Lock()
	defer f.m.Unlock()

	_, ok := f.s[c]
	if ok {
		delete(f.s, c)
		close(c)
	}
}

func (f *SseSplit) Publish(e string) {
	f.m.Lock()
	defer f.m.Unlock()

	if e == "" {
		return
	}

	for c := range f.s {
		select {
		case c <- e:
		default:
			delete(f.s, c)
			close(c)
		}
	}
}
