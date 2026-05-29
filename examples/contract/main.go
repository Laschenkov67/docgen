package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/laschenkov67/docgen"
)

func main() {
	eng := docgen.New()
	if err := eng.LoadTemplateFile("contract", "examples/contract/templates/contract.docx"); err != nil {
		log.Fatal(err)
	}

	data := map[string]any{
		"Number":    "12/2025",
		"Date":      time.Now(),
		"Customer":  map[string]any{"Name": "ООО «Краевед-Тур»", "INN": "1234567890"},
		"Performer": map[string]any{"Name": "ИП Петров П.П.", "INN": "500100732259"},
		"Amount":    150000.00,
	}

	f, _ := os.Create("contract.docx")
	defer f.Close()
	if err := eng.Render(context.Background(), "contract", data, docgen.FormatDOCX, f); err != nil {
		log.Fatal(err)
	}
	log.Println("ok: contract.docx")
}
