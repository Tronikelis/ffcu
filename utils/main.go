package utils

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/shirou/gopsutil/v3/process"
)

func KillProcess(name string) error {
	processes, err := process.Processes()
	if err != nil {
		return err
	}

	killed := false

	for _, p := range processes {
		n, err := p.Name()
		if err != nil {
			continue
		}

		if strings.Contains(strings.ToLower(n), strings.ToLower(name)) {
			if err := p.Kill(); err != nil {
				continue
			}

			killed = true
		}
	}

	if killed {
		return nil
	}

	return fmt.Errorf("can't find %q", name)
}

func DownloadBytes(u string) ([]byte, error) {
	resp, err := http.Get(u)
	if err != nil {
		return []byte{}, nil
	}

	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, nil
	}

	return bytes, nil
}
