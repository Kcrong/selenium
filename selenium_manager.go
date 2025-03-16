package selenigo

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// SeleniumManager is a wrapper for getting information from the Selenium Manager binaries
type SeleniumManager struct {
	binaryPath string
}

// LogEntry represents a log entry from the Selenium Manager
type LogEntry struct {
	Level   string `json:"level"`
	Message string `json:"message"`
}

// ManagerOutput represents the output from the Selenium Manager
type ManagerOutput struct {
	Logs   []LogEntry        `json:"logs"`
	Result map[string]string `json:"result"`
}

// NewSeleniumManager creates a new SeleniumManager instance
func NewSeleniumManager() *SeleniumManager {
	return &SeleniumManager{}
}

// BinaryPaths determines the locations of the requested assets
func (sm *SeleniumManager) BinaryPaths(args []string) (map[string]string, error) {
	binary, err := sm.getBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to get selenium manager binary: %v", err)
	}

	// Prepare command arguments
	cmdArgs := append([]string{binary}, args...)
	cmdArgs = append(cmdArgs,
		"--language-binding", "go",
		"--output", "json",
	)

	// Run the command
	result, err := sm.run(cmdArgs)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// getBinary determines the path of the correct Selenium Manager binary
func (sm *SeleniumManager) getBinary() (string, error) {
	if sm.binaryPath != "" {
		return sm.binaryPath, nil
	}

	// Check environment variable first
	if envPath := os.Getenv("SE_MANAGER_PATH"); envPath != "" {
		if _, err := os.Stat(envPath); err == nil {
			sm.binaryPath = envPath
			return envPath, nil
		}
	}

	// Get the binary name based on the platform
	binaryName := sm.getBinaryName()
	if binaryName == "" {
		return "", fmt.Errorf("unsupported platform/architecture combination: %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	// Try to find the binary in the same directory as the current file
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %v", err)
	}

	binaryPath := filepath.Join(filepath.Dir(exePath), "selenium-manager", binaryName)
	if _, err := os.Stat(binaryPath); err == nil {
		sm.binaryPath = binaryPath
		return binaryPath, nil
	}

	return "", fmt.Errorf("unable to obtain working Selenium Manager binary at: %s", binaryPath)
}

// getBinaryName returns the binary name based on the current platform
func (sm *SeleniumManager) getBinaryName() string {
	switch runtime.GOOS {
	case "windows":
		return "selenium-manager.exe"
	case "darwin":
		return "selenium-manager"
	case "linux", "freebsd", "openbsd":
		if runtime.GOARCH == "amd64" {
			return "selenium-manager"
		}
	}
	return ""
}

// run executes the Selenium Manager Binary
func (sm *SeleniumManager) run(args []string) (map[string]string, error) {
	cmd := exec.Command(args[0], args[1:]...)

	// Execute command and capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to execute selenium manager: %v\nOutput: %s", err, string(output))
	}

	// Parse JSON output
	var managerOutput ManagerOutput
	if err := json.Unmarshal(output, &managerOutput); err != nil {
		return nil, fmt.Errorf("failed to parse selenium manager output: %v", err)
	}

	// Process logs
	sm.processLogs(managerOutput.Logs)

	if cmd.ProcessState.ExitCode() != 0 {
		return nil, fmt.Errorf("selenium manager exited with code %d", cmd.ProcessState.ExitCode())
	}

	return managerOutput.Result, nil
}

// processLogs processes the log entries from the Selenium Manager
func (sm *SeleniumManager) processLogs(logs []LogEntry) {
	for _, log := range logs {
		switch log.Level {
		case "WARN":
			fmt.Printf("Warning: %s\n", log.Message)
		case "DEBUG", "INFO":
			fmt.Printf("Debug: %s\n", log.Message)
		}
	}
}
