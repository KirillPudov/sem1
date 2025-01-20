package main

import (
	"archive/zip"
	"net/http"
	"os"
	"strconv"
)

func giveData(res http.ResponseWriter, req *http.Request, db *DB) error {

	file_wr := createZip(res, req, db)

	res.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote("data.zip"))
	res.Header().Set("Content-Type", "application/octet-stream")

	http.ServeFile(res, req, file_wr)

	clearFiles(file_wr)

	return nil
}

func createZip(res http.ResponseWriter, req *http.Request, db *DB) string {
	archive, err := os.Create("data.zip")

	w := zip.NewWriter(archive)

	csv_file, err := w.Create("data.csv")

	defer archive.Close()

	err = db.exportToCSV(csv_file)

	if err != nil {
		http.Error(res, err.Error(), 500)
		return ""
	}

	w.Close()

	return archive.Name()

}

func clearFiles(file string) error {
	return os.Remove(file)
}
