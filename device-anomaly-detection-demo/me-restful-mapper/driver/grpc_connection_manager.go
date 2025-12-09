package driver

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
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
		if conn.GetState().String() != "SHUTDOWN" && conn.GetState().String() != "TRANSIENT_FAILURE" {
			return conn, nil
		}
		// Remove the broken connection
		klog.Warningf("Connection to %s is in state %s, creating a new one", address, conn.GetState().String())
		delete(m.connections, address)
	}

	// Check if the port is open before attempting gRPC connection
	dialer := net.Dialer{Timeout: 2 * time.Second}
	tcpConn, err := dialer.Dial("tcp", address)
	if err != nil {
		err := fmt.Errorf("port check failed for %s: %v", address, err)
		return nil, err
	}
	tcpConn.Close()

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
