package mysqldump

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
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

		log.Printf("skippedTables: %#v", skippedTables)
		for _, t := range skippedTables {
			dumpDataCmd.Args = append(dumpDataCmd.Args, fmt.Sprintf("--ignore-table=%v.%v", dbName, t))
		}

		log.Printf("args: %#v", dumpDataCmd.Args)
		dumpDataCmd.Stdout = gz
		dumpDataCmd.Stderr = errWriter
		dumpDataCmd.Run()
	}()

	defer pr.Close()
	io.Copy(outWriter, pr)

	return nil
}
