package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/laschenkov67/docgen"
)

func main() {
	tplPath := flag.String("t", "", "путь к шаблону (.docx или .pdf.tmpl)")
	dataPath := flag.String("d", "", "путь к JSON с данными")
	outPath := flag.String("o", "out", "путь к результату")
	format := flag.String("f", "", "формат: docx|pdf (если пусто — по расширению шаблона)")
	flag.Parse()

	if *tplPath == "" || *dataPath == "" {
		flag.Usage()
		os.Exit(2)
	}

	raw, err := os.ReadFile(*dataPath)
	must(err)
	var data any
	must(json.Unmarshal(raw, &data))

	f := docgen.Format(*format)
	if f == "" {
		switch strings.ToLower(filepath.Ext(*tplPath)) {
		case ".docx":
			f = docgen.FormatDOCX
		default:
			f = docgen.FormatPDF
		}
	}

	eng := docgen.New()
	must(eng.LoadTemplateFile("main", *tplPath))

	out, err := os.Create(*outPath)
	must(err)
	defer out.Close()
	must(eng.Render(context.Background(), "main", data, f, out))
	log.Printf("written: %s", *outPath)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
