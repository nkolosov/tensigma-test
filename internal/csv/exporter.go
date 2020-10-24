package csv

import (
	"fmt"
	"sync"

	"github.com/nkolosov/tendigma-test/internal/datasource"
)

type Exporter struct {
	ds    *datasource.Products
	files chan DownloadCSVResult
	wg    sync.WaitGroup
}

func NewExporter(ds *datasource.Products, files chan DownloadCSVResult) *Exporter {
	return &Exporter{
		ds:    ds,
		files: files,
	}
}

func (e *Exporter) RunWorkers(n int) {
	e.wg.Add(n)

	for i := 0; i < n; i++ {
		go e.run()
	}
}

func (e *Exporter) run() {
	defer e.wg.Done()

	for _ = range e.files {
		//todo: read file and insert data to MongoDB
		for i := 0; i < 10; i++ {
			e.ds.Update(datasource.NewProductModel(fmt.Sprintf("test-%d", i), uint64(i)))
		}
	}
}

func (e *Exporter) Close() {
	//channel must be closed in downloader previously
	e.wg.Wait()
}
