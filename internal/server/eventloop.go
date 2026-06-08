//go:build linux

package server

import (
	"fmt"
	"syscall"
	"time"

	"garnet/internal/logger"
	"garnet/internal/store"
)

func (s *Server) startEventLoop() error {
	epfd, err := syscall.EpollCreate1(0)
	if err != nil {
		return fmt.Errorf("epoll_create1 error: %v", err)
	}
	s.epfd = epfd
	defer syscall.Close(epfd)

	var event syscall.EpollEvent
	event.Events = syscall.EPOLLIN
	event.Fd = int32(s.fd)
	if err := syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, s.fd, &event); err != nil {
		return fmt.Errorf("epoll_ctl add server error: %v", err)
	}

	events := make([]syscall.EpollEvent, s.cfg.MaxClients)
	logger.Logger.Printf("epoll event loop started")

	lastExpireRun := time.Now()

	for {
		n, err := syscall.EpollWait(epfd, events, 100) // 100ms timeout
		if err != nil {
			if err == syscall.EINTR {
				continue
			}
			logger.Logger.Printf("epoll_wait error: %v", err)
			return err
		}

		for i := 0; i < n; i++ {
			fd := int(events[i].Fd)

			if fd == s.fd {
				s.accept()
			} else {
				conn := s.connections[fd]
				if conn != nil {
					if events[i].Events&syscall.EPOLLIN != 0 {
						s.read(conn)
					}
				}
			}
		}

		// Run active expiration if 100ms has passed
		if time.Since(lastExpireRun) >= 100*time.Millisecond {
			deleted := store.ActiveExpire()
			if deleted > 0 {
				logger.Logger.Printf("active expiration removed %d keys", deleted)
			}
			lastExpireRun = time.Now()
		}
	}
}

func (s *Server) accept() {
	nfd, _, err := syscall.Accept(s.fd)
	if err != nil {
		if err == syscall.EAGAIN || err == syscall.EWOULDBLOCK {
			return
		}
		logger.Logger.Printf("accept error: %v", err)
		return
	}

	if err := syscall.SetNonblock(nfd, true); err != nil {
		logger.Logger.Printf("setnonblock error on client: %v", err)
		syscall.Close(nfd)
		return
	}

	var event syscall.EpollEvent
	event.Events = syscall.EPOLLIN
	event.Fd = int32(nfd)
	if err := syscall.EpollCtl(s.epfd, syscall.EPOLL_CTL_ADD, nfd, &event); err != nil {
		logger.Logger.Printf("epoll_ctl add client error: %v", err)
		syscall.Close(nfd)
		return
	}

	s.connections[nfd] = NewConnection(nfd)
	logger.Logger.Printf("client connected: fd %d", nfd)
}

func (s *Server) read(conn *Connection) {
	err := conn.Read()
	if err != nil {
		s.closeConnection(conn.fd)
	}
}

func (s *Server) closeConnection(fd int) {
	logger.Logger.Printf("client disconnected: fd %d", fd)
	syscall.Close(fd)
	delete(s.connections, fd)
}
