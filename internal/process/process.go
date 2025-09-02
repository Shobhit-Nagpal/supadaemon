package process

import (
	"context"
	"errors"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/Shobhit-Nagpal/supadaemon/internal/data"
	"github.com/Shobhit-Nagpal/supadaemon/internal/utils"
)

type Process struct {
	processData data.Model
}

func New(data data.Model) *Process {
	return &Process{
		processData: data,
	}
}

func (p *Process) Spawn(ctx context.Context) {
	if isRunning() {
		return
	}

	ticker := time.NewTicker(p.processData.Interval())
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			url := p.processData.Url()
			go makeRequest(url)
		case <-ctx.Done():
			return
		}
	}
}

func makeRequest(url string) error {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 && resp.StatusCode < 600 {
		return errors.New("Request failed")
	}

	return nil
}

func isRunning() bool {
	pid, err := utils.ReadPID()
	if err != nil {
		return false
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	// Send signal 0 to test existence
	err = process.Signal(syscall.Signal(0))
	return err == nil
}
