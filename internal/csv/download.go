package csv

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"google.golang.org/grpc/grpclog"

	"io"
	"net/http"
	"os"
	"sync"
)

const (
	defaultBufferSize = 100
)

type DownloadCSVTask struct {
	url string
}

func NewDownloadCSVTask(url string) DownloadCSVTask {
	return DownloadCSVTask{
		url: url,
	}
}

type DownloadCSVResult struct {
	url          string
	filename     string
	isSuccess    bool
	errorMessage string
}

type Downloader struct {
	tasks   chan DownloadCSVTask
	results chan DownloadCSVResult

	wg                sync.WaitGroup
	downloadDirectory string
}

func NewDownloader(downloadDirectory string, tasks chan DownloadCSVTask, results chan DownloadCSVResult) *Downloader {
	return &Downloader{
		tasks:             tasks,
		results:           results,
		downloadDirectory: downloadDirectory,
		wg:                sync.WaitGroup{},
	}
}

func (d *Downloader) runWorkers(n int) {
	d.wg.Add(n)

	for i := 0; i < n; i++ {
		go d.run()
	}
}

func (d *Downloader) run() {
	defer d.wg.Done()

	var err error
	var filename string
	var result DownloadCSVResult
	var errorMessage string

	for task := range d.tasks {
		filename, err = downloadFile(d.downloadDirectory, task.url)
		if err != nil {
			grpclog.Warningf("can't download CSV file with error: %+v", err)
		}

		errorMessage = ""
		if err != nil {
			errorMessage = err.Error()
		}

		result = DownloadCSVResult{
			url:          task.url,
			filename:     filename,
			isSuccess:    err == nil,
			errorMessage: errorMessage,
		}

		d.results <- result
	}
}

func (d *Downloader) close() {
	d.wg.Wait()
}

func downloadFile(filePath string, url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", errors.Wrapf(err, "can't download file\n")
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	info, err := os.Stat(filePath)
	if err != nil {
		return "", errors.Wrapf(err, "can't get file info %s\n", filePath)
	}

	if !info.IsDir() {
		return "", fmt.Errorf("directory %s is not exists\n", filePath)
	}

	fileName := fmt.Sprintf("%s/%s.csv", filePath, uuid.New().String())

	out, err := os.Create(fileName)
	if err != nil {
		return "", errors.Wrapf(err, "can't create file %s\n", fileName)
	}

	defer func() {
		_ = out.Close()
	}()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", errors.Wrapf(err, "can't copy file %s to workdir", fileName)
	}

	return fileName, nil
}
