package cmd

import (
	"fmt"
	"log"

	"github.com/codenotary/immudb/embedded/sql"
	"github.com/codenotary/immudb/embedded/store"
	"github.com/sony/sonyflake"
)

func getFileId() uint64 {
	// Sonyflake Id
	var st sonyflake.Settings
	st.CheckMachineID = fakeMachineID
	sf := sonyflake.NewSonyflake(st)
	if sf == nil {
		log.Fatal("New Sonyflake failed!")
	}

	id, err := sf.NextID()
	if err != nil {
		log.Fatal("NextID failed!")
	}
	return id
}

type TTAsset struct {
	ttid     uint64
	hash     string
	filename string
	fileext  string
	url      string
}

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type DB struct {
	store  *store.ImmuStore
	engine *sql.Engine
	sqltx  *sql.SQLTx
}

func (db *DB) Open() {
	var err error
	// immudb embedded
	// create/open immudb store at specified path
	db.store, err = store.Open("data", store.DefaultOptions())
	handleErr(err)

	// initialize sql engine (specify a key-prefix to isolate generated kv entries)
	db.engine, err = sql.NewEngine(db.store, sql.DefaultOptions().WithPrefix([]byte("sql")))
	handleErr(err)

	_, _, err = db.engine.Exec("CREATE DATABASE IF NOT EXISTS tiktag;", nil, nil)
	handleErr(err)

	// set the database to use in the context of the ongoing sql tx
	_, _, err = db.engine.Exec("USE DATABASE tiktag;", nil, nil)
	handleErr(err)

	// a sql tx is created and carried over next statements
	db.sqltx, _, err = db.engine.Exec("BEGIN TRANSACTION;", nil, nil)
	handleErr(err)

	// creates a table
	_, _, err = db.engine.Exec(`
		CREATE TABLE IF NOT EXISTS TTAsset (
			ttid INTEGER,
			hash VARCHAR,
			filename VARCHAR,
			fileext VARCHAR,
			url VARCHAR,
			PRIMARY KEY ttid
		);`, nil, db.sqltx)
	handleErr(err)
}

func (db *DB) Close() {
	// ensure tx is closed (it won't affect committed tx)
	db.engine.Exec("ROLLBACK;", nil, db.sqltx)

	db.store.Close()
}

func (db *DB) Exec(sql string) {
	_, _, err := db.engine.Exec(sql, nil, db.sqltx)
	handleErr(err)

	// commit ongoing transaction
	_, _, err = db.engine.Exec("COMMIT;", nil, db.sqltx)
	handleErr(err)
}

func saveAsset(ttasset TTAsset) {
	db := new(DB)
	db.Open()
	defer db.Close()

	// insert some rows
	SQL := fmt.Sprintf(`
		INSERT INTO TTAsset (
			ttid,
			hash,
			filename,
			fileext,
			url
		) VALUES 
			(%d, '%s', '%s', '%s', '%s')
		;`, ttasset.ttid, ttasset.hash, ttasset.filename, ttasset.fileext, ttasset.url)
	// fmt.Println(SQL)
	db.Exec(SQL)
}

func searchAsset(fHash string) string {
	db := new(DB)
	db.Open()
	defer db.Close()

	// query data including ongoing and unconfirmed changes
	SQL := fmt.Sprintf(`
			SELECT url
			FROM ttasset
			WHERE hash = '%s';
	`, fHash)
	rowReader, err := db.engine.Query(SQL, map[string]interface{}{"value": 100}, db.sqltx)
	handleErr(err)

	// ensure row reader is closed
	defer rowReader.Close()

	// selected columns can be read from the rowReader
	cols, err := rowReader.Columns()
	handleErr(err)

	url := ""
	for {
		// iterate over result set
		row, err := rowReader.Read()
		if err == sql.ErrNoMoreRows {
			break
		}
		handleErr(err)

		// each row contains values for the selected columns
		url = fmt.Sprintf("%s", row.ValuesBySelector[cols[0].Selector()].Value())
		// log.Printf("row: %v\n", url)
	}

	// close row reader
	rowReader.Close()

	return url
}
