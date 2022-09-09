package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jschaf/pggen"
)

var (
	schemaPath = flag.String("schema", "./db/schema", "path to schema file")
	outputDir  = flag.String("output", "./internal/domain/chat/storage/postgres/", "path to output dir")
	queryFile  = flag.String("query", "./db/queries.sql", "path to queries file for pggen")
)

func main() {
	flag.Parse()
	connString := fmt.Sprintf("user=${PG_USER} password=${PG_PASSWORD} dbname=${PG_DATABASE}")
	connString = os.ExpandEnv(connString)
	err := pggen.Generate(pggen.GenerateOptions{
		SchemaFiles: []string{*schemaPath},
		ConnString:  connString,
		QueryFiles:  []string{*queryFile},
		OutputDir:   *outputDir,
		Language:    pggen.LangGo,
		TypeOverrides: map[string]string{
			"integer": "int64",
			"bigint":  "int64",
			"text":    "github.com/jackc/pgtype.Text",
		},
	})

	log.Fatal(err)
}
