package main

import (
	"fmt"
	"github.com/d3m0k1d/hyprd/pkg/ipc"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "hyprd",
	Short: "hyprd",
	Long:  `hyprd`,
}

func Init() {

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Init()
	Execute()
	listner := &ipc.Listener{}
	eventsChan, err := listner.Start()
	if err != nil {
		os.Exit(1)
	}
	for {
		select {
		case eventsChan := <-eventsChan:
			fmt.Println(eventsChan)
		}
	}
}
