package service

import (
	"encoding/json"
	"github-actions-manager/internal/core"
	"os"
	"path/filepath"
	"strings"
)

type runnerConfig struct {
	AgentName string `json:"agentName"`
	GitHubURL string `json:"gitHubUrl"`
}

func DiscoverRunners(basePath string) ([]core.Runner, error) {
	var runners []core.Runner

	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		return runners, nil
	}

	entries, err := os.ReadDir(basePath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		installPath := filepath.Join(basePath, entry.Name())
		configPath := filepath.Join(installPath, ".runner")

		if _, err := os.Stat(configPath); err == nil {
			data, err := os.ReadFile(configPath)
			if err != nil {
				continue
			}

			data = removeBOM(data)

			var config runnerConfig
			if err := json.Unmarshal(data, &config); err != nil {
				continue
			}

			parts := strings.Split(strings.TrimRight(config.GitHubURL, "/"), "/")
			projectName := parts[len(parts)-1]

			runners = append(runners, core.Runner{
				RunnerName:  config.AgentName,
				ProjectName: projectName,
				ProjectURL:  config.GitHubURL,
				InstallPath: installPath,
				Status:      core.StatusStopped,
				IsInstalled: false,
			})
		}
	}

	return runners, nil
}

func removeBOM(data []byte) []byte {
	if len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		return data[3:]
	}
	return data
}
