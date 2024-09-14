package db

import (
	"compress/gzip"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os/exec"

	"github.com/e10k/sshdbd/connections"
	"github.com/gliderlabs/ssh"
)

func Dump(s ssh.Session, conn *connections.Connection, dbName string, skippedTables []string, outWriter io.Writer, errWriter io.Writer) error {
	sessionId := s.Context().SessionID()[:10]

	pr, pw := io.Pipe()

	gz := gzip.NewWriter(pw)

	go func() {
		defer pw.Close()
		defer gz.Close()

		dumpSchemaCmd := exec.Command(
			"mysqldump",
			"--single-transaction",
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
		err := dumpSchemaCmd.Run()

		if err != nil {
			log.Printf("[%s] dumping interrupted: %v\n", sessionId, err)
			log.Printf("[%s] cleaning up (making sure process PID %v is gone)\n", sessionId, dumpSchemaCmd.Process.Pid)
			err2 := dumpSchemaCmd.Process.Kill()
			if err2 != nil {
				log.Printf("[%s] error killing process: %v\n", sessionId, err2)
			}
			return
		}

		dumpDataCmd := exec.Command(
			"mysqldump",
			"--single-transaction",
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
		err = dumpDataCmd.Run()

		if err != nil {
			log.Printf("[%s] dumping interrupted: %v\n", sessionId, err)
			log.Printf("[%s] cleaning up (making sure process PID %v is gone)\n", sessionId, dumpDataCmd.Process.Pid)
			err2 := dumpDataCmd.Process.Kill()
			if err2 != nil {
				log.Printf("[%s] error killing process: %v\n", sessionId, err2)
			}
			return
		}
	}()

	defer pr.Close()
	io.Copy(outWriter, pr)

	log.Printf("[%s] done\n", sessionId)
	errWriter.Write([]byte("\n🏁 Finished dumping.\n"))

	return nil
}

func GetDatabases(conn *connections.Connection) ([]string, error) {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/?charset=utf8mb4&parseTime=True&loc=Local", conn.Username, conn.Password, conn.Host, conn.Port)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed connecting to the database server: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed pinging the database server: %v\n", err)
	}

	rows, err := db.Query("SHOW DATABASES;")
	if err != nil {
		return nil, fmt.Errorf("failed listing the databases: %v\n", err)
	}
	defer rows.Close()

	var databases []string
	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			return nil, fmt.Errorf("failed fetching the databases: %v\n", err)
		}
		databases = append(databases, dbName)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error fetching the databases: %v\n", err)
	}

	return databases, nil
}
