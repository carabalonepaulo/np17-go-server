package listener

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type Event interface {
	Handle(l *Listener)
}

type ClientConnected struct{ id int }

func (c ClientConnected) Handle(l *Listener) {
	if l.OnClientConnected != nil {
		l.OnClientConnected(c.id)
	}
}

type ClientDisconnected struct{ id int }

func (c ClientDisconnected) Handle(l *Listener) {
	if l.OnClientDisconnected != nil {
		l.OnClientDisconnected(c.id)
	}
}

type MessageReceived struct {
	id      int
	message string
}

func (c MessageReceived) Handle(l *Listener) {
	if l.OnMessageReceived != nil {
		l.OnMessageReceived(c.id, c.message)
	}
}

type Client struct {
	id            int
	conn          net.Conn
	writer        chan string
	totalSent     int
	totalReceived int
}

type Listener struct {
	running  bool
	listener net.Listener
	events   chan Event
	clients  []*Client

	OnClientConnected    func(int)
	OnClientDisconnected func(int)
	OnMessageReceived    func(int, string)
}

func NewListener(host string, port int, capacity int) (*Listener, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, err
	}

	self := &Listener{
		listener: listener,
		running:  true,
		events:   make(chan Event),
		clients:  make([]*Client, capacity),
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, syscall.SIGINT)

	go self.captureCtrlC(sc)
	go self.beginAccept()

	return self, nil
}

func (s Listener) Running() bool {
	return s.running
}

// essa função traz, de forma não bloqueante, os eventos para a main thread
func (s *Listener) DispatchEvents() {
	select {
	case e := <-s.events:
		e.Handle(s)
	default:
	}
}

func (s *Listener) TotalSent(id int) int {
	if s.clients[id] == nil {
		return 0
	}
	return s.clients[id].totalSent
}

func (s *Listener) TotalReceived(id int) int {
	if s.clients[id] == nil {
		return 0
	}
	return s.clients[id].totalReceived
}

func (s *Listener) SendTo(id int, message string) bool {
	if s.clients[id] == nil {
		return false
	}
	s.clients[id].writer <- message
	return true
}

func (s *Listener) SendToMany(message string, filter func(int) bool) {
	for i := 0; i < len(s.clients); i++ {
		if s.clients[i] != nil && filter(i) {
			s.clients[i].writer <- message
		}
	}
}

func (s *Listener) SendToAll(message string) {
	for i := 0; i < len(s.clients); i++ {
		if s.clients[i] != nil {
			s.clients[i].writer <- message
		}
	}
}

func (s *Listener) Kick(id int) {
	s.clients[id].conn.Close()
}

func (s *Listener) KickAll() {
	for i := 0; i < len(s.clients); i++ {
		if s.clients[i] != nil {
			s.clients[i].conn.Close()
		}
	}
}

func (s *Listener) Close() {
	s.running = false
	s.listener.Close()
}

func (s *Listener) captureCtrlC(signals chan os.Signal) {
	<-signals
	s.Close()
}

func (s *Listener) beginAccept() {
	for s.running {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		c := &Client{id: 0, conn: conn, writer: make(chan string)}
		c.id = s.findFreeSlot()

		if c.id == -1 {
			// TODO: já está cheio
			panic("not implemented yet")
		}

		s.clients[c.id] = c
		fmt.Printf("[Accept] Slot[%d] == nil // %t\n", c.id, (s.clients[c.id] == nil))
		s.events <- ClientConnected{id: c.id}

		// enviar e receber são operações bloqueantas
		// logo é preciso separa-las, cada uma com sua goroutine
		go s.beginReceive(c)
		go s.beginSend(c)
	}
}

func (s *Listener) beginReceive(c *Client) {
	r := bufio.NewReader(c.conn)
	for s.running {
		message, err := r.ReadString('\n')
		if err != nil {
			s.events <- ClientDisconnected{id: c.id}
			s.clients[c.id] = nil
			return
		}
		c.totalReceived += int(len(message))
		s.events <- MessageReceived{id: c.id, message: message}
	}
}

func (s *Listener) beginSend(c *Client) {
	fail := func(err error) {
		log.Printf("Error while sending data to `%d`: %v", c.id, err)
		s.Kick(c.id) // TODO: ???
	}

	w := bufio.NewWriter(c.conn)
	for s.running {
		message := <-c.writer

		_, err := w.WriteString(message)
		if err != nil {
			fail(err)
		}
		c.totalSent += int(len(message))

		// se a mensagem não tiver um separador... evita recriar a string
		if message[len(message)-1] != '\n' {
			if w.WriteByte('\n') != nil {
				fail(err)
			}
			c.totalSent += 1
		}
	}
}

func (s *Listener) findFreeSlot() int {
	for i := 0; i < len(s.clients); i++ {
		if s.clients[i] == nil {
			return i
		}
	}
	return -1
}
