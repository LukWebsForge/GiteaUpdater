package main

import (
	"errors"
	"fmt"
	"github.com/mitchellh/ioprogress"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const downloadServer = "https://dl.gitea.io/gitea/"

func DownloadGitea(version string) (path string, err error) {

	// Creating the destination file
	file, err := os.Create("gitea-" + version)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Creating the url by the version and the running system
	url := downloadServer + version + "/gitea-" + version + "-" + runtime.GOOS + "-" + runtime.GOARCH

	// Downloading the content
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// If the file doesn't exists we return an error
	if resp.StatusCode == 404 {
		return "", errors.New("File " + resp.Request.URL.String() + " could not be found (404)!")
	}

	// Showing the progress
	progress := downloadProgress(resp.Body, resp.ContentLength)

	// Writing it to the file
	_, err = io.Copy(file, progress)
	if err != nil {
		return "", err
	}

	// Getting the full path of the file
	path, err = filepath.Abs(file.Name())
	if err != nil {
		return "", err
	}

	// Returning the full path
	return path, nil
}

func downloadProgress(reader io.ReadCloser, size int64) (progress io.Reader) {

	duration, _ := time.ParseDuration("20ms")

	drawFunc := ioprogress.DrawTerminalf(os.Stdout, func(progress, total int64) string {
		// Dividing to get to MegaBytes
		downloaded := float64(progress) / 1024.0 / 1024.0
		size := float64(total) / 1024.0 / 1024.0

		// Display the values
		return fmt.Sprintf("Downloading: %.2f MB/%.2f MB", downloaded, size)
	})

	progress = &ioprogress.Reader{
		Reader:       reader,
		Size:         size,
		DrawInterval: duration,
		DrawFunc:     drawFunc,
	}

	return progress
}
