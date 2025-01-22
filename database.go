package main

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"math"

	"github.com/gocarina/gocsv"
)

const (
	layoutISO = "2006-01-02"
)

func InitDb(db_host string, db_port int, db_user string, db_pass string, db_database string) (*DB, error) {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", db_host, db_port, db_user, db_pass, db_database)

	db_conn, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return nil, err
	}

	return &DB{
		conn: db_conn,
	}, err
}

func (db *DB) loadCSVtoDB(file io.ReadCloser) (*Responce, error) {

	csv_reader := csv.NewReader(file)
	csv_reader.Comma = ','
	csv_reader.LazyQuotes = true

	var entries []Entity
	err := gocsv.UnmarshalCSV(csv_reader, &entries)
	if err != nil {
		return nil, err
	}

	var total_items int
	var total_price float64
	var total_categories int

	fail := func(err error) (*Responce, error) {
		return nil, fmt.Errorf("%v", err)
	}

	ctx := context.Background()

	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		return fail(err)
	}

	defer tx.Rollback()

	for _, entry := range entries {
		// parse_date, _ := time.Parse(layoutISO, entry.Create_date)
		// date := parse_date.Format(layoutISO)
		_, err := tx.ExecContext(ctx, `insert into prices(name, category, price, create_date) values ($1, $2, $3, $4);`, entry.Name, entry.Category, entry.Price, entry.Create_date)
		if err != nil {
			return fail(err)
		}
	}

	if err = tx.QueryRowContext(ctx, "SELECT count(*) from prices;").Scan(&total_items); err != nil {
		if err == sql.ErrNoRows {
			return fail(fmt.Errorf("Total items error"))
		}
		return fail(err)
	}

	if err = tx.QueryRowContext(ctx, "SELECT sum(price) from prices;").Scan(&total_price); err != nil {
		if err == sql.ErrNoRows {
			return fail(fmt.Errorf("Total price error"))
		}
		return fail(err)
	}

	if err = tx.QueryRowContext(ctx, "SELECT COUNT(DISTINCT(category)) from prices;").Scan(&total_categories); err != nil {
		if err == sql.ErrNoRows {
			return fail(fmt.Errorf("Category error"))
		}
		return fail(err)
	}

	err = tx.Commit()
	if err != nil {
		return fail(err)
	}

	return &Responce{
		Total_items:      total_items,
		Total_categories: total_categories,
		Total_price:      math.Round(total_price*100) / 100,
	}, nil
}

func (db *DB) exportToCSV(file io.Writer) error {

	fail := func(err error) error {
		return fmt.Errorf("%v", err)
	}

	all, err := db.conn.Query(`SELECT * FROM prices;`)
	defer all.Close()

	if err != nil {
		return err
	}

	var csv_entity []Entity

	for all.Next() {
		var csv_str Entity
		if err := all.Scan(&csv_str.Id, &csv_str.Name, &csv_str.Category, &csv_str.Price, &csv_str.Create_date); err == sql.ErrNoRows {
			return fail(err)
		}
		csv_entity = append(csv_entity, csv_str)
	}

	csv_wr := csv.NewWriter(file)

	err = gocsv.MarshalCSV(csv_entity, csv_wr)
	if err != nil {
		return err
	}

	return nil
}
