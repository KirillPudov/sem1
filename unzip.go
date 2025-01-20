package main

import (
	"archive/zip"
	"bytes"
	"fmt"
)

func unzipAndInsertToDB(body []byte, db_conn *DB) (*Responce, error) {

	zipRead, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))

	if err != nil {
		return nil, err
	}

	var ansv *Responce

	for _, file := range zipRead.File {

		defer handlePanic()

		if file.FileInfo().IsDir() {
			continue
		}

		defer handlePanic()

		read_file, err := file.Open()

		if err != nil {
			return nil, err
		}

		ans, err := db_conn.loadCSVtoDB(read_file)
		if err != nil {
			return nil, err
		}

		ansv = ans
	}

	return ansv, nil

}

func handlePanic() {
	if r := recover(); r != nil {
		fmt.Println("Recovered, err:\n", r)
	}
}
