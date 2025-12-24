package main

import (
	"fmt"
	"github.com/d3m0k1d/hyprd/pkg/config"
	"github.com/d3m0k1d/hyprd/pkg/ipc"
	"github.com/d3m0k1d/hyprd/pkg/logger"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
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
	logger := logger.New(false)
	Init()
	Execute()
	cfg, err := config.LoadConfig()

	if err != nil {
		logger.Error("Failed to load config", "err", err)
		os.Exit(1)
	}

	triggerIndex := make(map[string][]config.Rule)
	for _, rule := range cfg.Rules {
		triggerIndex[rule.Trigger] = append(triggerIndex[rule.Trigger], rule)
	}
	logger.Info("Loaded rules", "count", len(cfg.Rules))

	listener := &ipc.Listener{}
	eventsChan, err := listener.Start()
	if err != nil {
		logger.Error("Failed to start listener", "err", err)
		os.Exit(1)
	}
	logger.Info("Event listener started")

	for {
		select {
		case event := <-eventsChan:
			logger.Info("Received event", "type", event.Event)

			if rules, exists := triggerIndex[event.Event]; exists {
				for _, rule := range rules {
					logger.Info("Executing rule", "name", rule.Name)
					executeActions(rule.Actions, event)
				}
			}
		}
	}
}

func executeActions(actions []string, event ipc.Event) {
	logger := logger.New(false)

	for _, action := range actions {
		logger.Info("Executing action", "action", action)

		cmd := exec.Command("sh", "-c", action)
		output, err := cmd.CombinedOutput()

		if err != nil {
			logger.Error("Action failed", "action", action, "err", err, "output", string(output))
		} else {
			logger.Debug("Action succeeded", "action", action)
		}
	}
}
