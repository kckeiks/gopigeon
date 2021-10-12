package gopigeon

import (
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"
)

const clientIDletters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var clientIDSet *idSet
var subscribers *Subscribers

type MQTTConn struct {
	Conn     net.Conn
	ClientID string
	Topics   []string
}

type Subscribers struct {
	mu          sync.Mutex
	subscribers map[string][]*MQTTConn
}

type idSet struct {
	mu  sync.Mutex
	set map[string]struct{}
}

func (s *Subscribers) addSubscriber(c *MQTTConn, topic string) {
	s.mu.Lock()
	s.subscribers[topic] = append(s.subscribers[topic], c)
	s.mu.Unlock()
}

func (s *Subscribers) removeSubscriber(c *MQTTConn, topic string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	subs, ok := s.subscribers[topic]
	if !ok {
		return
	}
	for i, sub := range subs {
		if c.ClientID == sub.ClientID {
			s.subscribers[topic] = append(s.subscribers[topic][:i], s.subscribers[topic][i+1:]...)
		}
	}
}

func (s *Subscribers) getSubscribers(topic string) ([]*MQTTConn, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	subs, ok := s.subscribers[topic]
	if ok {
		return subs, nil
	}
	return nil, UnknownTopicError
}

func NewClientID() string {
	rand.Seed(time.Now().UnixNano())
	var b strings.Builder
	b.Grow(MaxClientIDLen)
	for i := 0; i < MaxClientIDLen; i++ {
		b.WriteByte(clientIDletters[rand.Intn(len(clientIDletters))])
	}
	return b.String()
}

func IsValidClientID(id string) bool {
	if len(id) > MaxClientIDLen || len(id) == 0 {
		return false
	}
	for _, r := range id {
		if (r < '0' || r > '9') && (r < 'A' || r > 'Z') && (r < 'a' || r > 'z') {
			return false
		}
	}
	return true
}

func (s *idSet) saveClientID(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.set[id] = struct{}{}
}

func (s *idSet) isClientIDUnique(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.set[id]
	return !ok
}

func (s *idSet) removeClientID(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.set, id)
}
