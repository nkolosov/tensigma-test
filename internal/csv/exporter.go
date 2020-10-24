package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
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
		return errors.Wrapf(err, "can't open the csv file %s", filename)
	}

	defer func() {
		err = csvFile.Close()
		if err != nil {
			grpclog.Errorf("can't close CSV file with error %#v", err)
		}
	}()

	var name string
	var price int

	r := csv.NewReader(csvFile)
	r.Comma = ';'
	r.FieldsPerRecord = 2
	r.TrimLeadingSpace = true

	for {
		// Read each record from csv
		record, err := r.Read()

		grpclog.Infof("read record %#v %#v\n", record, err)

		if err == io.EOF {
			break
		}

		if err != nil {
			return errors.Wrapf(err, "can't read data from CSV %s\n", filename)
		}

		if len(record) != 2 {
			return fmt.Errorf("invalid CSV row %#v", record)
		}

		name = record[0]
		price, _ = strconv.Atoi(record[1])

		model := datasource.NewProductModel(name, uint64(price))

		grpclog.Infof("prepare model %+v\n", model)

		e.ds.Update(model)
	}

	return nil
}

func (e *Exporter) close() {
	//channel must be closed in downloader previously
	e.wg.Wait()
}
