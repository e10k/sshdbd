package mysqldump

import (
	"compress/gzip"
	"fmt"
	"io"
	"os/exec"

	"github.com/e10k/dbdl/config"
)

func Dump(conn *config.Connection, outWriter io.Writer, errWriter io.Writer) error {
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
			conn.Dbname,
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
			conn.Dbname,
			"--ignore-table=sakila.film",
		)
		dumpDataCmd.Stdout = gz
		dumpDataCmd.Stderr = errWriter
		dumpDataCmd.Run()
	}()

	defer pr.Close()
	io.Copy(outWriter, pr)

	return nil
}
