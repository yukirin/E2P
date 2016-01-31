package main

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
)

// EnExport is en-export tag
type EnExport struct {
	Note []Note `xml:"note"`
}

// Note is note tag
type Note struct {
	Title  string `xml:"title"`
	Source string `xml:"note-attributes>source-url"`
}

func convert(in io.Reader, out io.Writer) error {
	e := EnExport{}
	if err := xml.NewDecoder(in).Decode(&e); err != nil {
		return err
	}

	// source-urlが空だとpocketに追加されないため
	dummy := "https://getpocket.com/"
	for i, v := range e.Note {
		if v.Source == "" {
			e.Note[i].Source = dummy
		}
	}

	t := template.Must(template.New("pocket").Parse(tpl))
	if err := t.Execute(out, e); err != nil {
		return err
	}
	return nil
}

func main() {
	if len(os.Args) != 3 {
		log.Fatal("引数が間違っています")
	}

	in, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	defer in.Close()

	out, err := os.Create(os.Args[2])
	if err != nil {
		fmt.Println(err)
		return
	}
	defer out.Close()

	if err := convert(in, out); err != nil {
		fmt.Println(err)
		return
	}
}

var tpl = `
<!DOCTYPE html>
<html>
  <head>
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
    <title>Instapaper: Export</title>
  </head>
  <body>
    <h1>Unread</h1>
    <ol>
	{{range $i, $n := .Note}}
	    <li><a href="{{$n.Source}}">{{$n.Title}}</a></li>
	{{end}}
    </ol>
  </body>
</html>
`
