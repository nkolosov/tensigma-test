package csv

import (
	"encoding/csv"
	"io"
	"os"
	"sync"

	"github.com/pkg/errors"
	"google.golang.org/grpc/grpclog"

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

func (e *Exporter) runWorkers(n int) {
	e.wg.Add(n)

	for i := 0; i < n; i++ {
		go e.run()
	}
}

func (e *Exporter) run() {
	defer e.wg.Done()

	var err error
	for file := range e.files {
		if !file.isSuccess {
			grpclog.Warningf("skip export CSV data after error %s", file.errorMessage)
			continue
		}

		grpclog.Infof("read export task %+v", file)

		err = e.exportFile(file.filename)
		if err != nil {
			grpclog.Infof("error on export file data %+v", err)
		}

		grpclog.Infof("completed export task %+v", file)
	}
}

func (e *Exporter) exportFile(filename string) error {
	csvFile, err := os.Open(filename)
	if err != nil {
		return errors.Wrapf(err, "can't open CSV file %s\n", filename)
	}

	defer func() {
		err = csvFile.Close()
		if err != nil {
			grpclog.Errorf("can't close CSV file with error %#v", err)
		}
	}()

	r := csv.NewReader(csvFile)
	r.Comma = ';'
	r.FieldsPerRecord = 2
	r.TrimLeadingSpace = true

	var product *datasource.Product

	for {
		// Read each record from csv
		columns, err := r.Read()

		grpclog.Infof("read record %#v %#v\n", columns, err)

		if err == io.EOF {
			break
		}

		if err != nil {
			return errors.Wrapf(err, "can't read data from CSV %s\n", filename)
		}

		product, err = datasource.CreateProductFromCSV(columns)
		e.ds.Update(product)
	}

	return nil
}

func (e *Exporter) close() {
	//channel must be closed in downloader previously
	e.wg.Wait()
}
