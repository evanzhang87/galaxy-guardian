package util

import (
	"errors"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var cachePath string

func init() {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	cachePath = path.Join(dir, "cache")
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		_ = os.MkdirAll(cachePath, os.ModePerm)
	}
}

func DownloadFromFileserver(url, version, fileName, agent string) (string, error) {
	if url == "" {
		return "", errors.New("url is empty")
	}

	rand.Seed(time.Now().UnixNano())
	var waitTime = rand.Intn(60) + 5
	if strings.Contains(url, ".asc") {
		waitTime = rand.Intn(5)
	}
	time.Sleep(time.Duration(waitTime) * time.Second)

	agentCache := path.Join(cachePath, agent)
	if _, err := os.Stat(agentCache); os.IsNotExist(err) {
		os.MkdirAll(agentCache, os.ModePerm)
	}

	versionCache := path.Join(agentCache, version)
	if _, err := os.Stat(versionCache); os.IsNotExist(err) {
		os.MkdirAll(versionCache, os.ModePerm)
	}

	AbfileName := path.Join(versionCache, fileName)
	if _, err := os.Stat(AbfileName); err == nil {
		os.Remove(AbfileName)
	}
	response, err := http.Get(url)
	if err != nil {
		return "", errors.New("Error while downloading " + url)
	}

	defer response.Body.Close()
	if response.StatusCode == 404 {
		return "", errors.New("404 ERROR can't find version to download from " + url)
	}

	output, err := os.Create(AbfileName)
	if err != nil {
		return "", errors.New("can't create the download file " + AbfileName)
	}

	defer output.Close()

	_, err = io.Copy(output, response.Body)
	if err != nil {
		return "", errors.New("Error while downloading " + url)
	}
	return AbfileName, nil
}

func DownloadPublicKey(url, path string) error {
	response, err := http.Get(url)
	if err != nil {
		return errors.New("Error while downloading " + url)
	}

	defer response.Body.Close()
	if response.StatusCode == 404 {
		return errors.New("404 ERROR can't find version to download from " + url)
	}

	output, err := os.Create(path)
	if err != nil {
		return errors.New("can't create the download file " + path)
	}

	defer output.Close()

	_, err = io.Copy(output, response.Body)
	if err != nil {
		return errors.New("Error while downloading " + url)
	}
	return nil
}

func GetModeCache(mode string) string {
	return path.Join(cachePath, "mode")
}
