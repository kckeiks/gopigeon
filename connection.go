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
	Conn       net.Conn
	ClientID   string
	Topics     []string // might be able to remove this
	KeepAlive  int
	MsgManager *MessageManager
}

type MessageManager struct {
	queued        []*Message
	inflight      []*Message
	InflightCount uint
	InflightMut   sync.Mutex
}

type Message struct {
	State MessageState
	Data  []byte
}

type Subscribers struct {
	mu          sync.Mutex
	subscribers map[string][]*MQTTConn
}

type idSet struct {
	mu  sync.Mutex
	set map[string]struct{}
}

func (m *MessageManager) enqueue(msg *Message) {
	m.queued = append(m.queued, msg)
}

func (m *MessageManager) dequeue() *Message {
	if len(m.queued) == 0 {
		return nil
	}
	msg := m.queued[0]
	m.queued = m.queued[1:]
	return msg
}

func (m *MessageManager) clear() {
	m.queued = nil
}

func (m *MessageManager) len() int {
	return len(m.queued)
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
