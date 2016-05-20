package main

import (
	"errors"
	"net"
	"time"
)

type StoppableListener struct {
	*net.TCPListener
	chstop chan struct{}
}

func New(l net.Listener) (*StoppableListener, error) {
	tcpL, ok := l.(*net.TCPListener)

	if !ok {
		return nil, errors.New("cannot wrap passed listener")
	}

	return &StoppableListener{tcpL, make(chan struct{})}, nil
}

var ErrStopped = errors.New("listener stopped")

func (sl *StoppableListener) Accept() (net.Conn, error) {

	for {
		sl.SetDeadline(time.Now().Add(time.Second))

		conn, err := sl.TCPListener.Accept()

		select {
		case <-sl.chstop:
			return nil, ErrStopped
		default:
		}

		if err != nil {
			netErr, ok := err.(net.Error)
			if ok && netErr.Timeout() && netErr.Temporary() {
				continue
			}
		}

		return conn, err
	}
}

func (sl *StoppableListener) Stop() {
	close(sl.chstop)
}
