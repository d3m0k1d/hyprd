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
	log := logger.New(false)
	Init()
	Execute()
	cfg, err := config.LoadConfig()

	if err != nil {
		log.Error("Failed to load config", "err", err)
		os.Exit(1)
	}

	triggerIndex := make(map[string][]config.Rule)
	for _, rule := range cfg.Rules {
		triggerIndex[rule.Trigger] = append(triggerIndex[rule.Trigger], rule)
	}
	log.Info("Loaded rules", "count", len(cfg.Rules))

	listener := &ipc.Listener{}
	eventsChan, err := listener.Start()
	if err != nil {
		log.Error("Failed to start listener", "err", err)
		os.Exit(1)
	}
	log.Info("Event listener started")

	for {
		select {
		case event := <-eventsChan:
			if rules, exists := triggerIndex[event.Event]; exists {
				for _, rule := range rules {
					log.Info("Executing rule", "name", rule.Name)
					executeActions(rule.Actions, event, *log)
				}
			}
		}
	}
}

func executeActions(actions []string, event ipc.Event, log logger.Logger) {
	log.Info("Starting actions execution", "count", len(actions))

	for i, action := range actions {
		log.Info("Executing action", "index", i, "action", action)

		cmd := exec.Command("sh", "-c", action)
		cmd.Env = os.Environ()
		cmd.Dir = os.Getenv("HOME")

		output, err := cmd.CombinedOutput()

		log.Info("Command output",
			"action", action,
			"output", string(output),
			"exit_code", cmd.ProcessState)

		if err != nil {
			log.Error("Action failed",
				"action", action,
				"err", err.Error(),
				"output", string(output),
				"exit_code", cmd.ProcessState)
		} else {
			log.Info("Action succeeded",
				"action", action,
				"output", string(output))
		}
	}

	log.Info("Finished actions execution")
}
