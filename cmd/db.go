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

func saveAsset(ttasset TTAsset) {
	// immudb embedded
	// create/open immudb store at specified path
	db, err := store.Open("data", store.DefaultOptions())
	handleErr(err)
	defer db.Close()

	// initialize sql engine (specify a key-prefix to isolate generated kv entries)
	engine, err := sql.NewEngine(db, sql.DefaultOptions().WithPrefix([]byte("sql")))
	handleErr(err)

	_, _, err = engine.Exec("CREATE DATABASE IF NOT EXISTS tiktag;", nil, nil)
	handleErr(err)

	// set the database to use in the context of the ongoing sql tx
	_, _, err = engine.Exec("USE DATABASE tiktag;", nil, nil)
	handleErr(err)

	// a sql tx is created and carried over next statements
	sqltx, _, err := engine.Exec("BEGIN TRANSACTION;", nil, nil)
	handleErr(err)

	// ensure tx is closed (it won't affect committed tx)
	defer engine.Exec("ROLLBACK;", nil, sqltx)

	// creates a table
	_, _, err = engine.Exec(`
		CREATE TABLE IF NOT EXISTS TTAsset (
			ttid INTEGER,
			hash VARCHAR,
			filename VARCHAR,
			fileext VARCHAR,
			url VARCHAR,
			PRIMARY KEY ttid
		);`, nil, sqltx)
	handleErr(err)

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
	_, _, err = engine.Exec(SQL, nil, sqltx)
	handleErr(err)

	// // query data including ongoing and unconfirmed changes
	// rowReader, err := engine.Query(`
	// 		SELECT id, date, creditaccount, debitaccount, amount, description
	// 		FROM journal
	// 		WHERE amount > @value;
	// `, map[string]interface{}{"value": 100}, sqltx)
	// handleErr(err)

	// // ensure row reader is closed
	// defer rowReader.Close()

	// // selected columns can be read from the rowReader
	// cols, err := rowReader.Columns()
	// handleErr(err)

	// for {
	// 	// iterate over result set
	// 	row, err := rowReader.Read()
	// 	if err == sql.ErrNoMoreRows {
	// 		break
	// 	}
	// 	handleErr(err)

	// 	// each row contains values for the selected columns
	// 	log.Printf("row: %v\n", row.ValuesBySelector[cols[0].Selector()].Value())
	// }

	// // close row reader
	// rowReader.Close()

	// commit ongoing transaction
	_, _, err = engine.Exec("COMMIT;", nil, sqltx)
	handleErr(err)
}
