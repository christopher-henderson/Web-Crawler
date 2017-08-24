package repository

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DIR = "/tmp/crawler"

const (
	dbFile = "/Users/chris/Documents/crawler/pornhub.sqlite3"
	driver = "sqlite3"
)

const Schema = `CREATE TABLE IF NOT EXISTS 'Document' (
	'id'	INTEGER NOT NULL UNIQUE,
	'url'	TEXT NOT NULL,
	'content'	TEXT NOT NULL,
	PRIMARY KEY('id')
);`

const Insert = `INSERT INTO Document (url, content) VALUES (?, ?)`

func init() {
	ExecuteTransactionalDDL(Schema)
}

func Save(name string, content []byte) {
	ExecuteTransactionalDDL(Insert, name, content)
}

// func Save(name string, content []byte) {
// 	URL, _ := url.Parse(name)
// 	fqdn := path.Join(DIR, URL.Hostname())
// 	fd, err := os.Create(fqdn)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	defer fd.Close()
// 	fd.Write(content)
// }

func ExecuteTransactionalDDL(query string, args ...interface{}) error {
	transaction, err := getDB().Begin()
	defer transaction.Commit()
	if err != nil {
		return err
	}
	stmt, err := transaction.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err := stmt.Exec(args...); err != nil {
		transaction.Rollback()
		return err
	}
	return nil
}

var getDB = func() func() *sql.DB {
	db, err := sql.Open(driver, dbFile)
	if err != nil {
		panic(err)
	}
	return func() *sql.DB {
		return db
	}
}()
