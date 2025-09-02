package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/Shobhit-Nagpal/supadaemon/internal/data"
)

const configDirName = ".supadaemon"
const configFileName = "supadaemon.json"
const pidFile = "/var/run/supadaemon.pid"
const servicePath = "/etc/systemd/system/supadaemon.service"

const defaultUrl = "https://supabase.com/"
const defaultInterval = 5 * time.Hour

var defaultConfig data.Model = data.NewModel(defaultUrl, defaultInterval)

func GetConfigDirectory() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := path.Join(homeDir, configDirName)

	return configDir, nil
}

func GetConfigData() (*data.Model, error) {
	configDir, err := GetConfigDirectory()
	if err != nil {
		return nil, err
	}

	configFilePath := buildConfigPath(configDir, configFileName)

	if !fileExists(configFilePath) {
		return nil, errors.New("Config file not found")
	}

	fileData, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	var processData data.Model
	err = json.Unmarshal(fileData, &processData)
	if err != nil {
		return nil, err
	}

	return &processData, nil
}

func WritePID() error {
	pid := getPID()
	return os.WriteFile(pidFile, []byte(strconv.Itoa(pid)), 0644)
}

func ReadPID() (int, error) {
	data, err := os.ReadFile(pidFile)
	if err != nil {
		return 0, err
	}

	pidStr := strings.TrimSpace(string(data))
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return 0, err
	}

	return pid, nil
}

func GetPIDPath() string {
	return pidFile
}

func SetupServiceFile() error {
	// Check if running as root
	if os.Geteuid() != 0 {
		return errors.New("service setup requires root privileges (run with sudo)")
	}

	// Get executable path
	execPath, err := os.Executable()
	if err != nil {
		return err
	}

	// Get current user info (the user who ran sudo)
	username := os.Getenv("SUDO_USER")
	if username == "" {
		currentUser, err := user.Current()
		if err != nil {
			return err
		}
		username = currentUser.Username
	}

	// Service file content
	serviceContent := fmt.Sprintf(`[Unit]
Description=Supabase Daemon
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=%s
Group=%s
ExecStart=%s start
ExecStop=%s stop
Restart=always
RestartSec=10

StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
`, username, username, execPath, execPath)

	// Write service file
	err = os.WriteFile(servicePath, []byte(serviceContent), 0644)
	if err != nil {
		return err
	}

	// Reload systemd
	cmd := exec.Command("systemctl", "daemon-reload")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to reload systemd: %v", err)
	}

	return nil
}

func SetupConfig() error {
	configDir, err := GetConfigDirectory()
	if err != nil {
		return err
	}

	if !dirExists(configDir) {
		err = os.Mkdir(configDir, 0755)
		if err != nil {
			return err
		}
	}

	configFile := buildConfigPath(configDir, configFileName)
	if !fileExists(configFile) {
		jsonData, err := json.Marshal(defaultConfig)
		if err != nil {
			return err
		}

		err = os.WriteFile(configFile, jsonData, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

func CleanupService() error {
	var errs []error

	// 1. Stop the service if running
	if err := stopService(); err != nil {
		errs = append(errs, fmt.Errorf("failed to stop service: %v", err))
	}

	// 2. Disable the service
	if err := disableService(); err != nil {
		errs = append(errs, fmt.Errorf("failed to disable service: %v", err))
	}

	// 3. Remove service file
	if fileExists(servicePath) {
		if err := os.Remove(servicePath); err != nil {
			errs = append(errs, fmt.Errorf("failed to remove service file: %v", err))
		}
	}

	// 4. Remove PID file if exists
	if fileExists(pidFile) {
		if err := os.Remove(pidFile); err != nil {
			errs = append(errs, fmt.Errorf("failed to remove PID file: %v", err))
		}
	}

	// 5. Reload systemd
	if err := reloadSystemd(); err != nil {
		errs = append(errs, fmt.Errorf("failed to reload systemd: %v", err))
	}

	// Return combined errors if any
	if len(errs) > 0 {
		return fmt.Errorf("cleanup encountered errors: %v", errs)
	}

	return nil
}

func CleanupConfig() error {
	configDir, err := GetConfigDirectory()
	if err != nil {
		return err
	}

	if dirExists(configDir) {
		return os.RemoveAll(configDir) // Removes directory and all contents
	}

	return nil
}

// Full cleanup - everything
func CleanupAll() error {
	// Check root privileges for service cleanup
	if os.Geteuid() != 0 {
		return errors.New("cleanup requires root privileges (run with sudo)")
	}

	var errs []error

	// Clean service
	if err := CleanupService(); err != nil {
		errs = append(errs, err)
	}

	// Clean config (runs as original user)
	if err := CleanupConfig(); err != nil {
		errs = append(errs, fmt.Errorf("failed to cleanup config: %v", err))
	}

	if len(errs) > 0 {
		return fmt.Errorf("cleanup errors: %v", errs)
	}

	return nil
}

func buildConfigPath(dir, file string) string {
	return path.Join(dir, file)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return info.IsDir()
}

func getPID() int {
	return os.Getpid()
}

func stopService() error {
	cmd := exec.Command("systemctl", "stop", "supadaemon")
	return cmd.Run()
}

func disableService() error {
	cmd := exec.Command("systemctl", "disable", "supadaemon")
	return cmd.Run()
}

func reloadSystemd() error {
	cmd := exec.Command("systemctl", "daemon-reload")
	return cmd.Run()
}
