package ipc

import (
	"bufio"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/d3m0k1d/hyprd/pkg/logger"
)

type Event struct {
	Event string
	Date  string
}

type Listener struct {
	conn net.Conn
	ch   chan Event
}

func (l *Listener) getSocketPath() (socketPath string) {
	logger := logger.New(false)
	runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
	hyprlandInstance := os.Getenv("HYPRLAND_INSTANCE_SIGNATURE")
	socketPath = filepath.Join(
		runtimeDir,
		"hypr",
		hyprlandInstance,
		".socket2.sock",
	)
	logger.Info("SocketPath finded", "path", socketPath)
	return socketPath
}

func (l *Listener) Start() (<-chan Event, error) {
	logger := logger.New(false)
	socketPath := l.getSocketPath()
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		logger.Error("Error: %v\n", "err", err)
		os.Exit(1)
	}
	logger.Info("Connect to socket success", "path", socketPath)
	l.conn = conn
	l.ch = make(chan Event, 10)

	go func() {
		scanner := bufio.NewScanner(l.conn)
		for scanner.Scan() {
			line := scanner.Text()
			event := Event{
				Event: line,
				Date:  time.Now().String(),
			}
			l.ch <- event
			logger.Info("Event received", "event", line)
		}
		if err := scanner.Err(); err != nil {
			logger.Error("Scanner error: %v\n", "err", err)
		}
	}()

	return l.ch, nil
}

func (l *Listener) Stop() {
	logger := logger.New(false)
	if l.conn != nil {
		l.conn.Close()
		l.conn = nil
		logger.Info("Connection closed")

	}
}
