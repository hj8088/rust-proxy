package cmd

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func doSyncFromRemote(filePath, remoteURL string) error {
	var err error
	if err = os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}

	req, _ := http.NewRequest("GET", remoteURL, nil)
	req.Header.Set("Connection", "close")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	io.Copy(f, resp.Body)

	return nil
}
