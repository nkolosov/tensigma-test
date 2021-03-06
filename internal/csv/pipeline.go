package csv

import (
	"errors"
	"sync"

	"github.com/nkolosov/tendigma-test/internal/config"
	"github.com/nkolosov/tendigma-test/internal/datasource"
)

type Pipeline struct {
	downloader *Downloader
	exporter   *Exporter

	downloaderTasks chan DownloadCSVTask
	exporterTasks   chan DownloadCSVResult

	isStopped bool
	mutex     sync.Mutex
}

func NewPipeline(cfg config.PipelineConfig, ds *datasource.Products) *Pipeline {
	downloaderTasks := make(chan DownloadCSVTask, defaultBufferSize)
	exporterTasks := make(chan DownloadCSVResult, defaultBufferSize)

	downloader := NewDownloader(cfg.DownloaderTempDirectory, downloaderTasks, exporterTasks)
	exporter := NewExporter(ds, exporterTasks)

	downloader.runWorkers(cfg.DownloaderWorkersCount)
	exporter.runWorkers(cfg.ExporterWorkersCount)

	return &Pipeline{
		downloader:      downloader,
		exporter:        exporter,
		downloaderTasks: downloaderTasks,
		exporterTasks:   exporterTasks,
	}
}

func (p *Pipeline) Handle(url string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.isStopped {
		return
	}

	p.downloaderTasks <- NewDownloadCSVTask(url)
}

func (p *Pipeline) Close() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.isStopped {
		return errors.New("pipeline already stopped\n")
	}

	p.isStopped = true

	close(p.downloaderTasks)
	p.downloader.close()

	close(p.exporterTasks)
	p.exporter.close()

	return nil
}
