// +build ignore

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var (
	mimeType = map[string]string{
		".js":    "text/javascript",
		".css":   "text/css",
		".html":  "text/html; charset=UTF-8",
		".ico":   "image/vnd.microsoft.icon",
		".png":   "image/png",
		".svg":   "image/svg+xml",
		".otf":   "font/otf",
		".woff":  "font/woff",
		".woff2": "font/woff2",
		".ttf":   "font/ttf",
		".eot":   "font/eot",
		".txt":   "text/plain",
	}
	conv             = map[string]interface{}{"conv": fmtByteSlice}
	endpointTemplate = template.Must(
		template.New("").Funcs(conv).Parse(`package front

// Code generated by go generate; DO NOT EDIT.

import (
    "net/http"

    "github.com/labstack/echo/v4"
)

func (e *Endpoint) RegisterRoutes(g *echo.Group) {
{{- range $name, $staticFile := . }}
    g.GET("{{ $name }}", func(c echo.Context) error {
        return c.Blob(http.StatusOK, "{{ $staticFile.Mime }}", []byte{ {{ conv $staticFile.Data }} })
    })
{{- end }}
}
`),
	)
)

const (
	endpointFile     = "endpoint.go"
	distFolder       = "node/dist/node"
	defaultIndexHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>Node</title>
  <base href="/">

  <meta name="viewport" content="width=device-width, initial-scale=1">
  <link rel="icon" type="image/x-icon" href="favicon.ico">
</head>
<body>
<div>
  <p> This is the default index, looks like you forget to generate the angular app before generating the golang endpoint.</p
</div>
</body>
</html>`
)

func fmtByteSlice(s []byte) string {
	builder := strings.Builder{}

	for _, v := range s {
		builder.WriteString(fmt.Sprintf("%d,", int(v)))
	}

	return builder.String()
}

type staticFile struct {
	Data []byte
	Mime string
}

func main() {
	staticFiles := make(map[string]staticFile)
	if _, err := os.Stat(distFolder); os.IsNotExist(err) {
		staticFiles["/index.html"] = staticFile{
			Data: []byte(defaultIndexHTML),
			Mime: mimeType[".html"],
		}
		log.Print("dist folder doesn't exist, default index.html will be used")
	} else {
		err := filepath.Walk(distFolder, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				// skip directories
				return nil
			}
			relativePath := strings.TrimPrefix(filepath.ToSlash(path), distFolder)
			data, err := ioutil.ReadFile(path)
			if err != nil {
				// If file not reading
				log.Printf("Error reading %s: %s", path, err)
				return err
			}
			// Add file name to map

			staticFiles[relativePath] = staticFile{
				Data: data,
				Mime: mimeType[filepath.Ext(relativePath)],
			}

			return nil
		})
		if err != nil {
			log.Fatal("Error walking through embed directory:", err)
		}
	}

	if index, ok := staticFiles["/index.html"]; ok {
		staticFiles["/*"] = index
	}

	// Create buffer
	builder := &bytes.Buffer{}

	// Execute template
	if err := endpointTemplate.Execute(builder, staticFiles); err != nil {
		log.Fatal("Error executing template", err)
	}

	// Formatting generated code
	data, err := format.Source(builder.Bytes())
	if err != nil {
		log.Fatal("Error formatting generated code", err)
	}

	if err := ioutil.WriteFile(endpointFile, data, os.ModePerm); err != nil {
		log.Fatal("Error writing endpoint file", err)
	}
}
