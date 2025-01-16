package main

import (
	"database/sql"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/gocarina/gocsv"
)

const (
	layoutISO = "2006-01-02"
)

func loadCSVtoDB(file string) (int, error) {

	// db_host, exist := os.LookupEnv("POSTGRES_HOST")
	// if exist != true {
	// 	return 0, errors.New("ENV variable POSTGRES_HOST is not exist")
	// }
	// db_port, exist := os.LookupEnv("POSTGRES_PORT")
	// if exist != true {
	// 	return 0, errors.New("ENV variable POSTGRES_PORT is not exist")
	// }
	// db_user, exist := os.LookupEnv("POSTGRES_USER")
	// if exist != true {
	// 	return 0, errors.New("ENV variable POSTGRES_USER is not exist")
	// }
	// db_pass, exist := os.LookupEnv("POSTGRES_PASSWORD")
	// if exist != true {
	// 	return 0, errors.New("ENV variable POSTGRES_PASSWORD is not exist")
	// }
	// db_database, exist := os.LookupEnv("POSTGRES_DB")
	// if exist != true {
	// 	return 0, errors.New("ENV variable POSTGRES_DB is not exist")
	// }

	db_host := "localhost"
	db_port := 5432
	db_user := "validator"
	db_pass := "val1dat0r"
	db_database := "project-sem-1"

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", db_host, db_port, db_user, db_pass, db_database)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return 0, err
	}

	file_entry, err := os.Open(file)
	if err != nil {
		return 0, err
	}

	var entries []Entity
	err = gocsv.UnmarshalFile(file_entry, &entries)
	if err != nil {
		return 0, err
	}

	var total_items int

	for _, entry := range entries {
		parse_date, _ := time.Parse(layoutISO, entry.Create_date)
		date := parse_date.Format(layoutISO)
		insert := `insert into prices(id, name, category, price, create_date) values ($1, $2, $3, $4, $5);`
		_, err := db.Exec(insert, entry.Id, entry.Name, entry.Category, entry.Price, date)
		if err != nil {
			return 0, err
		}
		total_items += 1
	}

	return total_items, nil
}

func getResponceData(total_items int) (*Responce, error) {

	// db_host, exist := os.LookupEnv("POSTGRES_HOST")
	// if exist != true {
	// 	return nil, errors.New("ENV variable POSTGRES_HOST is not exist")
	// }
	// db_port, exist := os.LookupEnv("POSTGRES_PORT")
	// if exist != true {
	// 	return nil, errors.New("ENV variable POSTGRES_PORT is not exist")
	// }
	// db_user, exist := os.LookupEnv("POSTGRES_USER")
	// if exist != true {
	// 	return nil, errors.New("ENV variable POSTGRES_USER is not exist")
	// }
	// db_pass, exist := os.LookupEnv("POSTGRES_PASSWORD")
	// if exist != true {
	// 	return nil, errors.New("ENV variable POSTGRES_PASSWORD is not exist")
	// }
	// db_database, exist := os.LookupEnv("POSTGRES_DB")
	// if exist != true {
	// 	return nil, errors.New("ENV variable POSTGRES_DB is not exist")
	// }

	db_host := "localhost"
	db_port := 5432
	db_user := "validator"
	db_pass := "val1dat0r"
	db_database := "project-sem-1"

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", db_host, db_port, db_user, db_pass, db_database)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return nil, err
	}

	categories, err := db.Query(`SELECT count("category") FROM prices;`)
	defer categories.Close()

	if err != nil {
		return nil, err
	}
	var total_categories []int

	for categories.Next() {
		var categorie int
		if err := categories.Scan(&categorie); err != nil {
			return nil, err
		}
		total_categories = append(total_categories, categorie)
	}

	prices, err := db.Query(`SELECT sum("price") FROM prices;`)
	defer prices.Close()

	var total_price []float64

	for prices.Next() {
		var price float64
		if err := prices.Scan(&price); err != nil {
			return nil, err
		}
		total_price = append(total_price, price)
	}

	return &Responce{
		Total_items:      total_items,
		Total_categories: total_categories[0],
		Total_price:      math.Round(total_price[0]*100) / 100,
	}, nil
}

func exportToCSV() error {
	// db_host, exist := os.LookupEnv("POSTGRES_HOST")
	// if exist != true {
	// 	return nil, errors.New("ENV variable POSTGRES_HOST is not exist")
	// }
	// db_port, exist := os.LookupEnv("POSTGRES_PORT")
	// if exist != true {
	// 	return nil, errors.New("ENV variable POSTGRES_PORT is not exist")
	// }
	// db_user, exist := os.LookupEnv("POSTGRES_USER")
	// if exist != true {
	// 	return nil, errors.New("ENV variable POSTGRES_USER is not exist")
	// }
	// db_pass, exist := os.LookupEnv("POSTGRES_PASSWORD")
	// if exist != true {
	// 	return nil, errors.New("ENV variable POSTGRES_PASSWORD is not exist")
	// }
	// db_database, exist := os.LookupEnv("POSTGRES_DB")
	// if exist != true {
	// 	return nil, errors.New("ENV variable POSTGRES_DB is not exist")
	// }
	db_host := "localhost"
	db_port := 5432
	db_user := "validator"
	db_pass := "val1dat0r"
	db_database := "project-sem-1"

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", db_host, db_port, db_user, db_pass, db_database)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return err
	}

	all, err := db.Query(`SELECT * FROM prices;`)
	defer all.Close()

	if err != nil {
		return err
	}

	var csv []Entity

	for all.Next() {
		var csv_str Entity
		if err := all.Scan(&csv_str.Id, &csv_str.Name, &csv_str.Category, &csv_str.Price, &csv_str.Create_date); err != nil {
			return err
		}
		csv = append(csv, csv_str)
	}

	dst, err := os.Create("data.csv")

	srcFile, err := os.OpenFile(dst.Name(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 777)

	err = gocsv.Marshal(csv, srcFile)
	if err != nil {
		return err
	}

	return nil
}
