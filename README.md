# sitetheme

Mapea el contenido editable (`veltylabs/site_content`) a la plantilla de aterrizaje (`github.com/tinywasm/layout/landing`), generando artefactos de SEO (`JSON-LD`) y resúmenes para IA (`llms.txt`).

## Inicio Rápido

### Instalación

```bash
go get github.com/veltylabs/sitetheme
```

### Funciones Principales

`sitetheme` exporta tres funciones puras:

```go
// Landing arma la página del sitio a partir de su contenido. domain es el
// dominio público del sitio (sin esquema), usado para construir los
// canonical URL; no vive en sitecontent.Content porque no es contenido
// editable por el cliente.
func Landing(c sitecontent.Content, domain string) *landing.Page

// JSONLD produce el documento schema.org del rubro declarado en c.SEO.
func JSONLD(c sitecontent.Content) string

// LLMsTxt produce el resumen en markdown que leen los motores generativos.
func LLMsTxt(c sitecontent.Content) string
```

## Ejemplo de Uso

```go
package main

import (
	"fmt"
	"github.com/veltylabs/site_content"
	"github.com/veltylabs/sitetheme"
)

func main() {
	content := sitecontent.Content{
		Brand: sitecontent.Brand{Name: "Mi Negocio"},
		Seo:   sitecontent.SEO{Description: "Bienvenido a Mi Negocio"},
	}

	page := sitetheme.Landing(content, "minegocio.cl")
	pages := page.RenderPages()

	for _, p := range pages {
		fmt.Printf("Ruta: %s, Título: %s\n", p.Path, p.Doc.Title)
	}

	jsonld := sitetheme.JSONLD(content)
	llmstxt := sitetheme.LLMsTxt(content)

	fmt.Println("JSON-LD:", jsonld)
	fmt.Println("llms.txt:", llmstxt)
}
```

## Arquitectura

Consulta la documentación detallada en [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md).
