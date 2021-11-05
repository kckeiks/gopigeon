package internal

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/kckeiks/gopigeon/mqttlib"
)

const ClientIDletters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var defaultKeepAliveTime = 120
var subscriberTable = &Subscribers{Subscribers: make(map[string][]*Client)}
var clientIDSet = &idSet{set: make(map[string]struct{})}

//
type Client struct {
	Conn net.Conn
	ID   string
	Subs []*Subscription // might be able to remove this
	//Subs  []Subscription
	KeepAlive  int
	MsgManager *MessageManager
}

type Subscription struct {
	QoS   uint
	Topic string
}

type MessageManager struct {
	queued        []*Message
	inflight      []*Message
	inflightCount uint
	inflightMut   sync.Mutex
}

type Message struct {
	State mqttlib.MessageState
	Data  []byte
}

type Subscribers struct {
	Mu          sync.Mutex
	Subscribers map[string][]*Client
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

func (s *Subscribers) AddSubscriber(c *Client, topic string) {
	s.Mu.Lock()
	s.Subscribers[topic] = append(s.Subscribers[topic], c)
	s.Mu.Unlock()
}

func (s *Subscribers) RemoveSubscriber(c *Client, topic string) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	subs, ok := s.Subscribers[topic]
	if !ok {
		return
	}
	for i, sub := range subs {
		if c.ID == sub.ID {
			s.Subscribers[topic] = append(s.Subscribers[topic][:i], s.Subscribers[topic][i+1:]...)
		}
	}
}

func (s *Subscribers) GetSubscribers(topic string) ([]*Client, error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	subs, ok := s.Subscribers[topic]
	if ok {
		return subs, nil
	}
	return nil, mqttlib.UnknownTopicError
}

func NewClientID() string {
	rand.Seed(time.Now().UnixNano())
	var b strings.Builder
	b.Grow(mqttlib.MaxClientIDLen)
	for i := 0; i < mqttlib.MaxClientIDLen; i++ {
		b.WriteByte(ClientIDletters[rand.Intn(len(ClientIDletters))])
	}
	return b.String()
}

func HandleClient(c *Client) error {
	defer HandleDisconnect(c)
	fh, err := mqttlib.ReadFixedHeader(c.Conn)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if fh.PktType != mqttlib.Connect {
		fmt.Println(mqttlib.ExpectingConnectPktError)
		return mqttlib.ExpectingConnectPktError
	}
	err = handleConnect(c, fh)
	if err != nil {
		fmt.Println(err)
		return err
	}
	for {
		c.resetReadDeadline()
		fh, err := mqttlib.ReadFixedHeader(c.Conn)
		if err != nil {
			fmt.Println(err)
			return err
		}
		switch fh.PktType {
		case mqttlib.Publish:
			err = handlePublish(c, fh)
		case mqttlib.Subscribe:
			err = handleSubscribe(c, fh)
		case mqttlib.Pingreq:
			err = handlePingreq(c, fh)
		case mqttlib.Connect:
			err = mqttlib.SecondConnectPktError
		case mqttlib.Disconnect:
			return nil
		default:
			fmt.Println("warning: unknonw packet")
		}
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
}

func IsValidClientID(id string) bool {
	if len(id) > mqttlib.MaxClientIDLen || len(id) == 0 {
		return false
	}
	for _, r := range id {
		if (r < '0' || r > '9') && (r < 'A' || r > 'Z') && (r < 'a' || r > 'z') {
			return false
		}
	}
	return true
}

func (s *idSet) SaveClientID(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.set[id] = struct{}{}
}

func (s *idSet) IsClientIDUnique(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.set[id]
	return !ok
}

func (s *idSet) RemoveClientID(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.set, id)
}

func (c *Client) resetReadDeadline() {
	// We set Deadline to one and a half the Client keep alive value.
	// This value could be user given or the server's default value.
	keepAlive := time.Second * time.Duration(c.KeepAlive+c.KeepAlive/2)
	if c.KeepAlive%2 == 1 {
		// add half a second
		keepAlive += time.Millisecond * time.Duration(500)
	}
	c.Conn.SetReadDeadline(time.Now().Add(keepAlive))
}
