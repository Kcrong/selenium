package common

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

// Service represents a driver service that manages a driver server process
type Service struct {
	// Path to the executable
	Path string
	// Port the service should run on
	Port int
	// Environment variables for the new process
	Env map[string]string
	// Process handle
	process *os.Process
	// Command handle
	cmd *exec.Cmd
	// Log output writer
	logOutput io.Writer
	// Driver path environment key
	driverPathEnvKey string
}

// ServiceOption is a function that configures a Service
type ServiceOption func(*Service)

// WithPort sets the port for the service
func WithPort(port int) ServiceOption {
	return func(s *Service) {
		s.Port = port
	}
}

// WithLogOutput sets the log output for the service
func WithLogOutput(output io.Writer) ServiceOption {
	return func(s *Service) {
		s.logOutput = output
	}
}

// WithEnv sets the environment variables for the service
func WithEnv(env map[string]string) ServiceOption {
	return func(s *Service) {
		s.Env = env
	}
}

// WithDriverPathEnvKey sets the environment key for the driver path
func WithDriverPathEnvKey(key string) ServiceOption {
	return func(s *Service) {
		s.driverPathEnvKey = key
	}
}

// NewService creates a new Service instance
func NewService(path string, options ...ServiceOption) *Service {
	s := &Service{
		Path: path,
		Port: 0, // Will be assigned a free port
		Env:  make(map[string]string),
	}

	// Apply options
	for _, option := range options {
		option(s)
	}

	return s
}

// Start starts the Service
func (s *Service) Start() error {
	if s.Path == "" {
		return fmt.Errorf("service path cannot be empty")
	}

	// If port is 0, find a free port
	if s.Port == 0 {
		port, err := getFreePort()
		if err != nil {
			return fmt.Errorf("failed to find free port: %v", err)
		}
		s.Port = port
	}

	// Prepare command
	args := s.CommandLineArgs()
	cmd := exec.Command(s.Path, args...)

	// Set up environment
	if len(s.Env) > 0 {
		env := os.Environ()
		for k, v := range s.Env {
			env = append(env, fmt.Sprintf("%s=%s", k, v))
		}
		cmd.Env = env
	}

	// Set up output
	if s.logOutput != nil {
		cmd.Stdout = s.logOutput
		cmd.Stderr = s.logOutput
	}

	// Start the process
	if err := cmd.Start(); err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("'%s' executable may have wrong permissions: %v", filepath.Base(s.Path), err)
		}
		return fmt.Errorf("failed to start process: %v", err)
	}

	s.cmd = cmd
	s.process = cmd.Process

	// Wait for the service to be connectable
	if err := s.waitUntilConnectable(); err != nil {
		_ = s.Stop()
		return err
	}

	return nil
}

// Stop stops the Service
func (s *Service) Stop() error {
	if s.process == nil {
		return nil
	}

	// Try to send shutdown command
	s.sendRemoteShutdownCommand()

	// Kill the process
	if err := s.process.Kill(); err != nil {
		return fmt.Errorf("failed to kill process: %v", err)
	}

	s.process = nil
	s.cmd = nil
	return nil
}

// CommandLineArgs returns the command line arguments for the service
// This should be overridden by specific driver services
func (s *Service) CommandLineArgs() []string {
	return []string{}
}

// URL returns the service URL
func (s *Service) URL() string {
	return fmt.Sprintf("http://localhost:%d", s.Port)
}

// IsRunning checks if the service is running
func (s *Service) IsRunning() bool {
	if s.process == nil {
		return false
	}

	// Check if process is still running
	if err := s.process.Signal(syscall.Signal(0)); err != nil {
		return false
	}

	return true
}

// waitUntilConnectable waits until the service is connectable
func (s *Service) waitUntilConnectable() error {
	const maxAttempts = 70
	const initialDelay = 10 * time.Millisecond
	const maxDelay = 500 * time.Millisecond

	for i := 0; i < maxAttempts; i++ {
		// Check if process is still running
		if !s.IsRunning() {
			if s.cmd != nil && s.cmd.ProcessState != nil {
				return fmt.Errorf("service process exited with status: %v", s.cmd.ProcessState.ExitCode())
			}
			return fmt.Errorf("service process unexpectedly exited")
		}

		// Try to connect
		if isPortConnectable(s.Port) {
			return nil
		}

		// Calculate delay with exponential backoff
		delay := initialDelay + time.Duration(i)*50*time.Millisecond
		if delay > maxDelay {
			delay = maxDelay
		}
		time.Sleep(delay)
	}

	return fmt.Errorf("failed to connect to the service after %d attempts", maxAttempts)
}

// sendRemoteShutdownCommand attempts to send a shutdown command to the service
func (s *Service) sendRemoteShutdownCommand() {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Try to send shutdown command, ignore errors
	_, _ = client.Get(fmt.Sprintf("%s/shutdown", s.URL()))
}

// isPortConnectable checks if a port is connectable
func isPortConnectable(port int) bool {
	addr := fmt.Sprintf("localhost:%d", port)
	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

// getFreePort finds a free port on the system
func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}

	return l.Addr().(*net.TCPAddr).Port, l.Close()
}
