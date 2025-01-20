package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

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

type DB struct {
	conn *sql.DB
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

	db_host := os.Getenv("DB_HOST")
	db_port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	db_user := os.Getenv("DB_USER")
	db_pass := os.Getenv("DB_PASSWORD")
	db_database := os.Getenv("DB_NAME")

	if err != nil {
		fmt.Errorf(err.Error())
		return
	}

	db_conn, err := InitDb(db_host, db_port, db_user, db_pass, db_database)

	if err != nil {
		fmt.Errorf(err.Error())
		return
	}

	if req.Method == http.MethodPost {
		uploadData(res, req, db_conn)
	} else if req.Method == http.MethodGet {
		giveData(res, req, db_conn)
	} else {
		http.Error(res, "Only GET or POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
}

func uploadData(res http.ResponseWriter, req *http.Request, db_conn *DB) {
	req.ParseMultipartForm(10 << 20)

	file, _, err := req.FormFile("file")

	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}

	defer file.Close()

	body, err := io.ReadAll(file)

	resp, err := unzipAndInsertToDB(body, db_conn)

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
