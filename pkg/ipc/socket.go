package ipc

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path/filepath"
)

func CreateConnect() {
	runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
	hyprlandInstance := os.Getenv("HYPRLAND_INSTANCE_SIGNATURE")
	socketPath := filepath.Join(
		runtimeDir,
		"hypr",
		hyprlandInstance,
		".socket2.sock",
	)

	fmt.Printf("Connecting to: %s\n", socketPath)

	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Connected!")

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("Event: %s\n", line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Scanner error: %v\n", err)
	}
}
