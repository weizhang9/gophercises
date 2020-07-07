package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"regexp"

	_ "github.com/lib/pq"
)

const (
	host   = "localhost"
	port   = "5432"
	user   = "Wei.Zhang"
	pass   = "RichardPooPoo"
	dbname = "gophercise_phone"
)

var phones = []string{
	"1234567890",
	"123 456 7891",
	"(123) 456 7892",
	"(123) 456-7893",
	"123-456-7894",
	"123-456-7890",
	"1234567892",
	"(123)456-7892",
}

func normalise(phone string) string {
	var buf bytes.Buffer
	for _, ch := range phone {
		if ch >= '0' && ch <= '9' {
			buf.WriteRune(ch)
		}
	}
	return buf.String()
}

func regexNormalise(phone string) string {
	re := regexp.MustCompile("[^0-9]")
	return re.ReplaceAllString(phone, "")
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable", host, port, user, pass)
	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	must(err, "failed to open db with new dbname")
	defer db.Close()
	must(db.Ping(), "failed to ping db")
	must(createPhoneNumberTable(db), "failed to create the table")
	for _, p := range phones {
		id, err := insertPhone(db, p)
		must(err, "failed to insert phone")
		must(normalisePhone(db, id), "failed to normalise phone")
	}
}

func must(err error, str string) {
	if err != nil {
		log.Fatalln(str, err)
	}
}

func createDB(db *sql.DB, name string) error {
	_, err := db.Exec("CREATE DATABASE " + name)
	if err != nil {
		return err
	}
	return nil
}

func resetDB(db *sql.DB, name string) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + name)
	if err != nil {
		return err
	}
	return createDB(db, name)
}

func createPhoneNumberTable(db *sql.DB) error {
	stmt := `
		CREATE TABLE IF NOT EXISTS phone_numbers (
			id SERIAL,
			value VARCHAR(255)
		)`
	_, err := db.Exec(stmt)
	return err
}

func insertPhone(db *sql.DB, phone string) (int, error) {
	stmt := `INSERT INTO phone_numbers(value) VALUES($1) RETURNING id`
	var id int
	err := db.QueryRow(stmt, phone).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func getPhone(db *sql.DB, id int) (string, error) {
	stmt := `SELECT * FROM phone_numbers WHERE id=$1`
	var number string
	err := db.QueryRow(stmt, id).Scan(&id, &number)
	if err != nil {
		return "", err
	}
	return number, nil
}

func normalisePhone(db *sql.DB, id int) error {
	number, err := getPhone(db, id)
	if err != nil {
		return err
	}
	newN := normalise(number)
	stmt := `UPDATE phone_numbers SET value=$1 WHERE id=$2`
	_, err = db.Exec(stmt, newN, id)
	if err != nil {
		return err
	}
	return nil
}
