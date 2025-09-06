package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/Shobhit-Nagpal/supadaemon/internal/data"
	"github.com/Shobhit-Nagpal/supadaemon/internal/errors"
)

const serviceName = "supadaemon"
const configDirName = ".supadaemon"
const configFileName = "supadaemon.json"
const pidFile = "/var/run/supadaemon.pid"
const servicePath = "/etc/systemd/system/supadaemon.service"

const defaultUrl = "https://supabase.com/"
const defaultInterval = 5 * time.Hour

var defaultConfig data.Model = data.NewModel(defaultUrl, defaultInterval)

func getConfigDirectory() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.NewFileError("Home directory", errors.ReadFile)
	}

	configDir := path.Join(homeDir, configDirName)

	return configDir, nil
}

func getConfigData() (*data.Model, error) {
	configDir, err := getConfigDirectory()
	if err != nil {
		return nil, err
	}

	configFilePath := buildConfigPath(configDir, configFileName)

	if !fileExists(configFilePath) {
		return nil, errors.NewFileError(configFilePath, errors.FileNotFound)
	}

	fileData, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, errors.NewFileError(configFilePath, errors.ReadFile)
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
	err := os.WriteFile(pidFile, []byte(strconv.Itoa(pid)), 0644)

	if err != nil {
		return errors.NewFileError(pidFile, errors.WriteFile)
	}

	return err
}

func ReadPID() (int, error) {
	data, err := os.ReadFile(pidFile)
	if err != nil {
		return 0, errors.NewFileError(pidFile, errors.ReadFile)
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
		return errors.NewServiceError(serviceName, errors.RootAccess)
	}

	// Get executable path
	execPath, err := os.Executable()
	if err != nil {
		return errors.NewSystemError(errors.ExecuteableLookup)
	}

	// Get current user info (the user who ran sudo)
	username := os.Getenv("SUDO_USER")
	if username == "" {
		currentUser, err := user.Current()
		if err != nil {
			return errors.NewSystemError(errors.UserLookup)
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
		return errors.NewFileError(servicePath, errors.WriteFile)
	}

	// Reload systemd
	err = reloadSystemd()
	if err != nil {
		return errors.NewSystemError(errors.ReloadDaemon)
	}

	return nil
}

func SetupConfig() error {
	configDir, err := getConfigDirectory()
	if err != nil {
		return err
	}

	if !dirExists(configDir) {
		err = os.Mkdir(configDir, 0755)
		if err != nil {
			return errors.NewFileError(configDir, errors.WriteFile)
		}
	}

	configFile := buildConfigPath(configDir, configFileName)
	if !fileExists(configFile) {
		jsonData, err := json.Marshal(defaultConfig)
		if err != nil {
			return errors.NewSystemError(errors.DataSerialization)
		}

		err = os.WriteFile(configFile, jsonData, 0644)
		if err != nil {
			return errors.NewFileError(configFile, errors.WriteFile)
		}
	}

	return nil
}

func SetupAll() error {
	err := SetupConfig()
	if err != nil {
		return err
	}

	err = SetupServiceFile()
	if err != nil {
		return err
	}

	return nil
}

func CleanupService() error {

	if err := stopService(); err != nil {
		return errors.NewServiceError(serviceName, errors.StopService)
	}

	if err := disableService(); err != nil {
		return errors.NewServiceError(serviceName, errors.DisableService)
	}

	if fileExists(servicePath) {
		if err := os.Remove(servicePath); err != nil {
			return errors.NewFileError(servicePath, errors.DeleteFile)
		}
	}

	// 4. Remove PID file if exists
	if fileExists(pidFile) {
		if err := os.Remove(pidFile); err != nil {
			return errors.NewFileError(pidFile, errors.DeleteFile)
		}
	}

	// 5. Reload systemd
	if err := reloadSystemd(); err != nil {
		return errors.NewSystemError(errors.ReloadDaemon)
	}

	return nil
}

func CleanupConfig() error {
	configDir, err := getConfigDirectory()
	if err != nil {
		return err
	}

	if dirExists(configDir) {
		return os.RemoveAll(configDir) // Removes directory and all contents
	}

	return nil
}

func CleanupAll() error {
	// Check root privileges for service cleanup
	if os.Geteuid() != 0 {
		return errors.NewServiceError(serviceName, errors.RootAccess)
	}

	// Clean service
	if err := CleanupService(); err != nil {
		return err
	}

	// Clean config (runs as original user)
	if err := CleanupConfig(); err != nil {
		return err
	}

	return nil
}

func EnableService() error {
	if os.Geteuid() != 0 {
		return errors.NewServiceError(serviceName, errors.RootAccess)
	}

	if isServiceEnabled() {
		return nil
	}

	err := enableService()
	if err != nil {
		return errors.NewServiceError(serviceName, errors.EnableService)
	}

	return nil
}

func StartService() error {
	if os.Geteuid() != 0 {
		return errors.NewServiceError(serviceName, errors.RootAccess)
	}

	if !isServiceEnabled() {
		err := EnableService()
		if err != nil {
			return err
		}
	}

	if isServiceRunning() {
		return nil
	}

	err := startService()
	if err != nil {
		errors.NewServiceError(serviceName, errors.StartService)
	}

	return nil
}

func StopService() error {
	if os.Geteuid() != 0 {
		return errors.NewServiceError(serviceName, errors.RootAccess)
	}

	if !isServiceRunning() {
		return nil
	}

	err := stopService()
	if err != nil {
		return errors.NewServiceError(serviceName, errors.StopService)
	}

	return nil
}

func DisableService() error {
	if !isServiceEnabled() {
		return nil
	}

	err := disableService()
	if err != nil {
		return errors.NewServiceError(serviceName, errors.DisableService)
	}

	return nil
}

func GetServiceStatus() error {
	status, err := getServiceStatus()
	if err != nil {
		return errors.NewServiceError(serviceName, errors.GetServiceStatus)
	}

	fmt.Print(status)

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

func startService() error {
	cmd := exec.Command("systemctl", "start", serviceName)
	return cmd.Run()
}

func enableService() error {
	cmd := exec.Command("systemctl", "enable", serviceName)
	return cmd.Run()
}

func stopService() error {
	cmd := exec.Command("systemctl", "stop", serviceName)
	return cmd.Run()
}

func disableService() error {
	cmd := exec.Command("systemctl", "disable", serviceName)
	return cmd.Run()
}

func getServiceStatus() (string, error) {
	cmd := exec.Command("systemctl", "status", serviceName)
	output, err := cmd.CombinedOutput()

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.ExitCode() <= 4 {
				return string(output), nil
			}
		}

		return "", err
	}

	return string(output), nil
}

func reloadSystemd() error {
	cmd := exec.Command("systemctl", "daemon-reload")
	return cmd.Run()
}

func isServiceEnabled() bool {
	cmd := exec.Command("systemctl", "is-active", serviceName)
	err := cmd.Run()

	if err != nil {
		return false
	}

	return true
}

func isServiceRunning() bool {
	cmd := exec.Command("systemctl", "is-active", serviceName)
	err := cmd.Run()

	if err != nil {
		return false
	}

	return true
}
