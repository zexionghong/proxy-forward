package socks_proxy_server

import (
	"fmt"
	"net"
	"proxy-forward/config"
	"proxy-forward/pkg/logging"
	"sync"
)

type Server struct {
	listener     *net.TCPListener
	mutex        sync.Mutex
	enableUPD    bool
	readTimeout  int
	writeTimeout int

	headRequest *RequestList
}

// NewServer return a new server.
func NewServer() *Server {
	fmt.Print(config.RuntimeViper.GetString("socks_proxy_server.port"))
	port := config.RuntimeViper.GetString("socks_proxy_server.port")
	if port == "" {
		port = ":3333"
	}
	addr, err := net.ResolveTCPAddr("tcp", port)
	if err != nil {
		logging.Log.Fatal(err.Error())
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		logging.Log.Fatal(err.Error())
	}

	return &Server{
		listener:     l,
		mutex:        sync.Mutex{},
		enableUPD:    config.RuntimeViper.GetBool("socks_proxy_server.udp_support"),
		readTimeout:  config.RuntimeViper.GetInt("socks_proxy_server.read_timeout"),
		writeTimeout: config.RuntimeViper.GetInt("socks_proxy_server.write_timeout"),
	}
}

// Serve Run
func (s *Server) Run() error {
	logging.Log.Info("start server on: %s", s.listener.Addr().String())
	var er error
	for {
		conn, err := s.listener.AcceptTCP()
		er = err
		if err != nil {
			break
		}
		go connHandle(conn, s)
	}
	s.listener = nil
	s.closeRequests()
	return er
}

func (s *Server) closeRequests() {
	for s.headRequest != nil {
		r := s.headRequest
		s.removeRequestList(r)
		_ = r.Data.Close()
	}
}
func (s *Server) insertRequestList(l *RequestList) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// fifo 队列
	if s.headRequest != nil {
		s.headRequest.Prev = l
		l.Next = s.headRequest
		s.headRequest = l
	} else {
		s.headRequest = l
	}
}
func (s *Server) removeRequestList(l *RequestList) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.headRequest == l {
		s.headRequest = l.Next
	}
	if l.Prev != nil {
		l.Prev.Next = l.Next
		l.Prev = nil
	}
	if l.Next != nil {
		l.Next.Prev = l.Prev
		l.Next = nil
	}
}
func connHandle(conn *net.TCPConn, s *Server) {
	r := &RequestList{
		Prev: nil, Next: nil, Data: Request{ClientConn: conn, server: s},
	}
	s.insertRequestList(r)
	r.Data.Process()
	s.removeRequestList(r)
	_ = r.Data.Close()
}
