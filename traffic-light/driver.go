package main

import (
	"fmt"
	"log"
	"os/exec"
)

const (
	UNKNOW = "unknow"
	ON     = "ON"
	OFF    = "OFF"
)

func Set(number int64, state string) error {
	switch state {
	case ON:
		return On(number)
	case OFF:
		return Off(number)
	}
	return fmt.Errorf("unsupport state %v", state)
}

func On(number int64) error {
	on := exec.Command("gpio", "write", fmt.Sprintf("%d", number), "1")
	_, err := on.CombinedOutput()
	return err
}

func Off(number int64) error {
	on := exec.Command("gpio", "write", fmt.Sprintf("%d", number), "0")
	_, err := on.CombinedOutput()
	return err
}

func State(number int64) (string, error) {
	on := exec.Command("gpio", "read", fmt.Sprintf("%d", number))
	state, err := on.CombinedOutput()
	return string(state), err
}

func SetOutput(number int64) {
	cmd := exec.Command("gpio", "mode", fmt.Sprintf("%d", number), "output")
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd run error %v", err)
	}
}
