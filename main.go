package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Entity struct {
	Id          int     `csv:"id",db:"id"`
	Name        string  `csv:"name",db:"name"`
	Category    string  `csv:"category",db:"category"`
	Price       float64 `csv:"price",db:"price"`
	Create_date string  `csv:"create_date",db:"create_date"`
}

type DBEntity struct {
	Id          int     `db:"id"`
	Name        string  `db:"name"`
	Category    string  `db:"category"`
	Price       float64 `db:"price"`
	Create_date string  `db:"create_date"`
}

type Responce struct {
	Total_items      int     `json:"total_items"`
	Total_categories int     `json:"total_categories"`
	Total_price      float64 `json:"total_price"`
}

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/prices", handleRequest)

	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		panic(err)
	}
}
func handleRequest(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		uploadData(res, req)
	} else if req.Method == http.MethodGet {
		giveData(res, req)
	} else {
		http.Error(res, "Only GET or POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
}

func uploadData(res http.ResponseWriter, req *http.Request) {
	req.ParseMultipartForm(10 << 20)

	file, handler, err := req.FormFile("file")

	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}

	defer file.Close()

	file_path := filepath.Join(os.TempDir(), handler.Filename)

	dst, err := os.Create(file_path)
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}

	defer dst.Close()

	io.Copy(dst, file)

	resp, err := unzipAndInsertToDB(file_path)

	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}

	responceJson, err := json.Marshal(resp)
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}

	res.Header().Set("content-type", "application/json")

	res.WriteHeader(http.StatusOK)

	res.Write(responceJson)

}
