package ipc

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"
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
	runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
	hyprlandInstance := os.Getenv("HYPRLAND_INSTANCE_SIGNATURE")
	socketPath = filepath.Join(
		runtimeDir,
		"hypr",
		hyprlandInstance,
		".socket2.sock",
	)

	return socketPath
}

func (l *Listener) Start() (<-chan Event, error) {
	socketPath := l.getSocketPath()
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

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
		}
		if err := scanner.Err(); err != nil {
			fmt.Printf("Scanner error: %v\n", err)
		}
	}()

	return l.ch, nil
}

func (l *Listener) Stop() {
	if l.conn != nil {
		l.conn.Close()
		l.conn = nil
	}
}

