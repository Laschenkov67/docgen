# docgen

[![CI](https://github.com/laschenkov67/docgen/actions/workflows/ci.yml/badge.svg)](...)
[![Go Reference](https://pkg.go.dev/badge/github.com/laschenkov67/docgen.svg)](https://pkg.go.dev/github.com/yourorg/docgen)

Универсальный генератор документов (DOCX, PDF) из шаблонов для Go-проектов: интернет-магазинов.

## Возможности

- DOCX-шаблоны с синтаксисом `text/template` ({{.Field}}, {{range}}, FuncMap).
- Корректная работа с «разорванными runs» Microsoft Word.
- PDF-генерация без внешних зависимостей (gofpdf, UTF-8).
- Pluggable `Renderer`-интерфейс — добавляйте свои форматы.
- Готовые форматтеры под русскую локаль: `rubles`, `money`, `date`, `dateLong`,
  `inn`, `kpp`, валидаторы ИНН.