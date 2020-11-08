package service

import (
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/domain"
	"sync"
)

type happinessDoorPubSub struct {
	mu     sync.RWMutex
	subs   map[int][]chan domain.HappinessDoorDto
	closed bool
}

func newPubsub() *happinessDoorPubSub {
	ps := &happinessDoorPubSub{}
	ps.subs = make(map[int][]chan domain.HappinessDoorDto)
	return ps
}

func (ps *happinessDoorPubSub) publish(msg domain.HappinessDoorDto) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	if ps.closed {
		return
	}

	for _, ch := range ps.subs[msg.Id] {
		ch <- msg
	}
}

func (ps *happinessDoorPubSub) subscribe(hdID int) <-chan domain.HappinessDoorDto {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan domain.HappinessDoorDto, 1)
	ps.subs[hdID] = append(ps.subs[hdID], ch)
	return ch
}

func (ps *happinessDoorPubSub) unsubscribe(hdID int, ch <-chan domain.HappinessDoorDto) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ind, found := ps.getSubIndex(hdID, ch)
	if found {
		close(ps.subs[hdID][ind])
		ps.subs[hdID] = append(ps.subs[hdID][:ind], ps.subs[hdID][ind:]...)
	}
}

func (ps *happinessDoorPubSub) close() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if !ps.closed {
		ps.closed = true
		for _, subs := range ps.subs {
			for _, ch := range subs {
				close(ch)
			}
		}
	}
}

func (ps *happinessDoorPubSub) getSubIndex(hdID int, ch <-chan domain.HappinessDoorDto) (int, bool) {
	for i, sub := range ps.subs[hdID] {
		if sub == ch {
			return i, true
		}
	}
	return 0, false
}
