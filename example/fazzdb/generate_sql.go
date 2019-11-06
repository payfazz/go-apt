package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/payfazz/go-apt/example/fazzdb/fazzdb_sample/migration"
	"github.com/payfazz/go-apt/pkg/fazzdb"
)

const (
	SQL_DIR           = "./sql"
	SQL_FORMATTER_CLI = "sql-formatter-cli"
)

var migrations = []fazzdb.MigrationVersion{
	migration.Version1,
	migration.Version2,
}

func main() {
	fmt.Println("generating sql migration files")

	if !isSQLFormatterCLIExist() {
		panic("sql-formatter-cli not found: please install sql-formatter-cli before continue; 'npm install -g sql-formatter-cli'")
	}

	queries := fazzdb.Raw(true, migrations...)

	err := os.MkdirAll(SQL_DIR, os.FileMode(0744))
	if nil != err {
		panic(err)
	}

	for i, v := range queries {
		generatedFile := fmt.Sprintf("%s/%d.sql", SQL_DIR, i+1)

		_, err = os.Stat(generatedFile)
		if !os.IsNotExist(err) {
			fmt.Println("file", generatedFile, "already exists, skip generating current migration")
			continue
		}

		err = ioutil.WriteFile(generatedFile, []byte(v), os.FileMode(0644))
		if nil != err {
			panic(err)
		}

		cmd := exec.Command(
			SQL_FORMATTER_CLI,
			"-i", generatedFile, "-o", generatedFile,
		)

		if err = cmd.Start(); nil != err {
			fmt.Println("ERROR: failed formatting file", generatedFile)
			continue
		}
	}

	fmt.Println("finish generate migration into sql files")
}

func isSQLFormatterCLIExist() bool {
	_, err := exec.LookPath(SQL_FORMATTER_CLI)
	return nil == err
}
