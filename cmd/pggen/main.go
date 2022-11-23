package main

import (
	"flag"
	"log"
	"os"

	"github.com/jschaf/pggen"
)

var (
	schemaPath = flag.String("schema", "./db/schema.sql", "path to schema file")
	outputDir  = flag.String("output", "", "path to output dir")
	queryFile  = flag.String("query", "./db/queries.sql", "path to queries file for pggen")
)

func main() {
	flag.Parse()
	if *outputDir == "" {
		log.Fatal("output dir must be not empty")
	}
	connString := os.ExpandEnv("user=${PG_USER} password=${PG_PASSWORD} dbname=${PG_DATABASE}")
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

	if err != nil {
		log.Fatal(err)
	}
	log.Print("Successfully")
}
