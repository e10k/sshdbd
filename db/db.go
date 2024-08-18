package db

import (
	"compress/gzip"
	"database/sql"
	"fmt"
	"io"
	"os/exec"

	"github.com/e10k/dbdl/config"
)

func Dump(conn *config.Connection, dbName string, skippedTables []string, outWriter io.Writer, errWriter io.Writer) error {
	pr, pw := io.Pipe()

	gz := gzip.NewWriter(pw)

	go func() {
		defer pw.Close()
		defer gz.Close()

		dumpSchemaCmd := exec.Command(
			"mysqldump",
			"--single-transaction",
			"--databases",
			"--no-data",
			"--skip-add-drop-table",
			"--verbose",
			"-h",
			conn.Host,
			fmt.Sprintf("-u%s", conn.Username),
			fmt.Sprintf("-p%s", conn.Password),
			dbName,
		)
		dumpSchemaCmd.Stdout = gz
		dumpSchemaCmd.Stderr = errWriter
		dumpSchemaCmd.Run()

		dumpDataCmd := exec.Command(
			"mysqldump",
			"--single-transaction",
			"--tz-utc=false",
			"--no-create-info",
			"--verbose",
			"-h",
			conn.Host,
			fmt.Sprintf("-u%s", conn.Username),
			fmt.Sprintf("-p%s", conn.Password),
			dbName,
		)

		for _, t := range skippedTables {
			dumpDataCmd.Args = append(dumpDataCmd.Args, fmt.Sprintf("--ignore-table=%v.%v", dbName, t))
		}

		dumpDataCmd.Stdout = gz
		dumpDataCmd.Stderr = errWriter
		dumpDataCmd.Run()
	}()

	defer pr.Close()
	io.Copy(outWriter, pr)

	return nil
}

func GetDatabases(conn *config.Connection) ([]string, error) {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/?charset=utf8mb4&parseTime=True&loc=Local", conn.Username, conn.Password, conn.Host, conn.Port)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed connecting to the database: %v", err)
	}
	defer db.Close()

	rows, err := db.Query("SHOW DATABASES;")
	if err != nil {
		return nil, fmt.Errorf("failed listing the databases: %v", err)
	}
	defer rows.Close()

	var databases []string
	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			return nil, fmt.Errorf("failed fetching the databases: %v", err)
		}
		databases = append(databases, dbName)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error fetching the databases: %v", err)
	}

	return databases, nil
}
