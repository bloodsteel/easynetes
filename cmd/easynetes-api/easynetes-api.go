package main

import (
	"github.com/bloodsteel/easynetes/cmd"
	"github.com/bloodsteel/easynetes/pkg/log"
)

func main() {
	EasynetesAPICmd := cmd.EasynetesAPI()
	if err := EasynetesAPICmd.Execute(); err != nil {
		log.Errorf("exectu error: %v", err)
	}
}
