package driver

import (
	"context"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"k8s.io/klog/v2"
)

type GRPCConnectionManager struct {
	connections map[string]*grpc.ClientConn
	mutex       sync.Mutex
}

func NewGRPCConnectionManager() *GRPCConnectionManager {
	return &GRPCConnectionManager{
		connections: make(map[string]*grpc.ClientConn),
	}
}

func (m *GRPCConnectionManager) GetConnection(address string) (*grpc.ClientConn, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Check if a connection already exists
	if conn, exists := m.connections[address]; exists {
		state := conn.GetState()
		if state == connectivity.Ready || state == connectivity.Idle {
			return conn, nil
		}
		// Close and remove connections that are not reusable
		klog.Warningf("Connection to %s is in state %s, closing and creating a new one", address, state)
		_ = conn.Close()
		delete(m.connections, address)
	}

	// Create a new gRPC connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		err := fmt.Errorf("failed to create gRPC connection to %s: %v", address, err)
		return nil, err
	}

	m.connections[address] = conn
	return conn, nil
}

func (m *GRPCConnectionManager) CloseAllConnections() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for address, conn := range m.connections {
		if err := conn.Close(); err != nil {
			klog.Errorf("Failed to close connection to %s, error: %v", address, err)
		} else {
			klog.Infof("Closed connection to %s", address)
		}
	}

	m.connections = make(map[string]*grpc.ClientConn)
}
