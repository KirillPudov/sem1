package main

import (
	"io"
	"net/http"
	"os"
	"strconv"

	"archive/zip"
)

func giveData(res http.ResponseWriter, req *http.Request) {
	err := exportToCSV()
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}

	err = createZip()
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}

	res.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote("data.zip"))
	res.Header().Set("Content-Type", "application/octet-stream")

	http.ServeFile(res, req, "data.zip")

}

func createZip() error {
	archive, err := os.Create("data.zip")

	defer archive.Close()

	zipWriter := zip.NewWriter(archive)

	csv_file, err := os.OpenFile("data.csv", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 777)

	defer zipWriter.Close()

	if err != nil {
		return err
	}

	defer csv_file.Close()

	zip_archive, err := zipWriter.Create("data.csv")

	if _, err := io.Copy(zip_archive, csv_file); err != nil {
		return err
	}

	return nil
}
