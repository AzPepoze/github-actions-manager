package ops

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

type DownloadProgress struct {
	Downloaded int64
	Total      int64
}

func DownloadRunner(url string, destination string, progress chan<- DownloadProgress) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download runner: %s", response.Status)
	}

	file, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer file.Close()

	totalSize := response.ContentLength
	var downloadedSize int64

	buffer := make([]byte, 32*1024)
	for {
		read, readErr := response.Body.Read(buffer)
		if read > 0 {
			_, writeErr := file.Write(buffer[:read])
			if writeErr != nil {
				return writeErr
			}
			downloadedSize += int64(read)
			if progress != nil {
				progress <- DownloadProgress{
					Downloaded: downloadedSize,
					Total:      totalSize,
				}
			}
		}
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return readErr
		}
	}

	return nil
}

func ExtractRunner(projectName string, archivePath string) (string, error) {
	installPath := fmt.Sprintf("./actions/%s", projectName)
	if err := os.MkdirAll(installPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	command := exec.Command("tar", "-xzf", archivePath, "-C", installPath)
	if err := command.Run(); err != nil {
		return "", fmt.Errorf("failed to extract runner: %w", err)
	}

	return installPath, nil
}
