package gopigeon

import (
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"
)

const clientIDletters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var defaultKeepAliveTime = 120
var subscribers = &Subscribers{subscribers: make(map[string][]*MQTTConn)}
var clientIDSet = &idSet{set: make(map[string]struct{})}

type MQTTConn struct {
	Conn      net.Conn
	ClientID  string
	Topics    []string // might ne anle to remove this
	KeepAlive int
	// queue map[string]MessageQueue
	// inflightMax
	// inflightCount
}

// db__message_write_inflight_out_single

type MessageQueue struct {
	queue       *queue
	lastMsgSent []byte
	state       uint
}

// Queue for chunks of data
type queue struct {
	data [][]byte
}

func (q *queue) enqueue(b []byte) {
	q.data = append(q.data, b)
}

func (q *queue) dequeue() []byte {
	if len(q.data) == 0 {
		return nil
	}
	b := q.data[0]
	q.data = q.data[1:]
	return b
}

func (q *queue) clear() {
	q.data = nil
}

func (q *queue) len() int {
	return len(q.data)
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
