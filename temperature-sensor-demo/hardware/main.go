package main

import (
	"log"

	"github.com/thinkgos/gomodbus"
)

func InitModbusSimulator(address string) {
	srv := modbus.NewTCPServer()
	srv.AddNodes(
		modbus.NewNodeRegister(1, 0, 1, 0, 0, 0,0,0,1),
	)
	defer srv.Close()
	if err := srv.ListenAndServe(address); err != nil {
		log.Fatalf("Failed to start modbus simulator: %v", err)
	}
}
func main() {
	InitModbusSimulator(":5502")
}
