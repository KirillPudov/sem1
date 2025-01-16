package main

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"fmt"
)

func unzipAndInsertToDB(file string) (*Responce, error) {

	files, err := zip.OpenReader(file)
	if err != nil {
		return nil, err
	}

	defer files.Close()

	var total_items int

	for _, file := range files.File {

		file_path := filepath.Join(os.TempDir(), file.Name)

		if !strings.HasPrefix(file_path, filepath.Clean(os.TempDir())+string(os.PathSeparator)) {
			return &Responce{}, errors.New("File path error")
		}

		defer handlePanic()

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(file_path, os.ModePerm); err != nil {
				panic(err)
			}
			continue
		}

		defer handlePanic()

		if err := os.MkdirAll(filepath.Dir(file_path), os.ModePerm); err != nil {
			panic(err)
			continue
		}

		dstFile, err := os.OpenFile(file_path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return nil, err
		}

		srcFile, err := file.Open()
		if err != nil {
			return nil, err
		}

		if _, err := io.Copy(dstFile, srcFile); err != nil {
			return nil, err
		}

		total_items, err = loadCSVtoDB(file_path); if err != nil {
			return nil, err
		}

		defer dstFile.Close()
		defer srcFile.Close()
	}

	return getResponceData(total_items)

}

func handlePanic() {
		if r := recover(); r != nil {
			fmt.Println("Recovered, err:\n", r)
		}
}
