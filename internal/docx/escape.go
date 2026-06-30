package docx

import (
	"reflect"
	"strings"
)

// maxEscapeDepth защищает от паники переполнения стека на циклических
// или патологически глубоких структурах данных.
const maxEscapeDepth = 64

// escapeTemplateData возвращает глубокую копию data, в которой все строковые
// значения экранированы под XML-текст (см. escapeXMLText). Без этого шага
// значения вроде `ООО «Ромашка & Ко»` или адрес с символом "<" ломают
// word/document.xml: text/template не делает автоэкранирование, в отличие
// от html/template, а его контекстный HTML5-парсер не годится для OOXML-разметки
// (там полно неизвестных тегов вроде <w:t>, самозакрывающихся <w:rPr/> и т. п.).
//
// Экранирование выполняется на уровне данных, а не результата рендера, чтобы
// не задеть статическую разметку самого шаблона (она доверенная — её редактирует
// автор шаблона). Используются числовые ссылки (&#38; и т. п.), а не именованные
// (&amp;), потому что они нечувствительны к регистру и переживают применение
// FuncMap-функций upper/lower к уже экранированной строке.
func escapeTemplateData(data any) any {
	if data == nil {
		return nil
	}
	return escapeValue(reflect.ValueOf(data), 0).Interface()
}

func escapeValue(v reflect.Value, depth int) reflect.Value {
	if !v.IsValid() || depth >= maxEscapeDepth {
		return v
	}

	switch v.Kind() {
	case reflect.Interface:
		return escapeInterface(v, depth)
	case reflect.Ptr:
		return escapePtr(v, depth)
	case reflect.String:
		return escapeString(v)
	case reflect.Slice:
		return escapeSlice(v, depth)
	case reflect.Array:
		return escapeArray(v, depth)
	case reflect.Map:
		return escapeMap(v, depth)
	case reflect.Struct:
		return escapeStruct(v, depth)
	default:
		return v
	}
}

func escapeInterface(v reflect.Value, depth int) reflect.Value {
	if v.IsNil() {
		return v
	}
	out := reflect.New(v.Type()).Elem()
	out.Set(escapeValue(v.Elem(), depth+1))
	return out
}

func escapePtr(v reflect.Value, depth int) reflect.Value {
	if v.IsNil() {
		return v
	}
	out := reflect.New(v.Type().Elem())
	out.Elem().Set(escapeValue(v.Elem(), depth+1))
	return out
}

func escapeString(v reflect.Value) reflect.Value {
	out := reflect.New(v.Type()).Elem()
	out.SetString(escapeXMLText(v.String()))
	return out
}

func escapeSlice(v reflect.Value, depth int) reflect.Value {
	if v.IsNil() {
		return v
	}
	out := reflect.MakeSlice(v.Type(), v.Len(), v.Len())
	for i := 0; i < v.Len(); i++ {
		out.Index(i).Set(escapeValue(v.Index(i), depth+1))
	}
	return out
}

func escapeArray(v reflect.Value, depth int) reflect.Value {
	out := reflect.New(v.Type()).Elem()
	for i := 0; i < v.Len(); i++ {
		out.Index(i).Set(escapeValue(v.Index(i), depth+1))
	}
	return out
}

func escapeMap(v reflect.Value, depth int) reflect.Value {
	if v.IsNil() {
		return v
	}
	out := reflect.MakeMapWithSize(v.Type(), v.Len())
	iter := v.MapRange()
	for iter.Next() {
		out.SetMapIndex(escapeValue(iter.Key(), depth+1), escapeValue(iter.Value(), depth+1))
	}
	return out
}

func escapeStruct(v reflect.Value, depth int) reflect.Value {
	out := reflect.New(v.Type()).Elem()
	out.Set(v) // сначала копируем как есть — сохраняет неэкспортируемые поля
	for i := 0; i < v.NumField(); i++ {
		if !v.Type().Field(i).IsExported() {
			continue
		}
		out.Field(i).Set(escapeValue(v.Field(i), depth+1))
	}
	return out
}

func escapeXMLText(s string) string {
	if !strings.ContainsAny(s, "&<>\"'") {
		return s
	}
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch r {
		case '&':
			b.WriteString("&#38;")
		case '<':
			b.WriteString("&#60;")
		case '>':
			b.WriteString("&#62;")
		case '"':
			b.WriteString("&#34;")
		case '\'':
			b.WriteString("&#39;")
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}
