---
PLAN: "feat: mapeo de contenido a la plantilla landing con JSON-LD y llms.txt"
TAG: v0.1.0
EXECUTOR: jules
REVIEWER: none
STATUS: review
SESSION: 2203076454838285662
PR: https://github.com/veltylabs/sitetheme/pull/1
---

> Este plan se despacha con el flujo CodeJob. Ver skill: agents-workflow.
>
> **PUERTA: no despachar hasta que `github.com/veltylabs/site_content` v0.1.0 esté
> publicado.** Este plan importa sus tipos; sin ellos no compila.

# Plan — `veltylabs/sitetheme` v0.1.0

Crear el **tema** de los sitios de cliente de Velty: la pieza que convierte el
contenido editable en una página de `tinywasm/layout/landing`, más los
artefactos de SEO/GEO que dependen del negocio.

## La propiedad que hay que preservar por encima de todo

Este código lo ejecutan **dos** consumidores:

- el **panel** (`veltylabs/misitio`), en el navegador, para previsualizar;
- el **repositorio de sitios** (`veltylabs/clientsites`), en CI, para construir.

Como es el mismo código, la vista previa **no puede desviarse** de lo que se
publica. Todo lo demás en este plan está subordinado a eso: cualquier cosa que
haga que el resultado dependa del entorno (una fecha, un aleatorio, una lectura
de red, un estado global) rompe la propiedad aunque compile.

---

## 0. Reglas de desarrollo — léelas completas antes de escribir código

Están en [`AGENTS.md`](../AGENTS.md). Lo que sigue es lo que no puedes dejar de
aplicar.

### 0.1 Este repositorio compila **para wasm**

El panel lo usa en el navegador. Por lo tanto:

- El archivo principal (`landing.go` y equivalentes) **no lleva** `//go:build !wasm`.
- Sólo CSS/SVG/JS/HTML pesado va en archivos con nombre de extensión y
  `//go:build !wasm`: `css.go`, `js.go`, `html.go`, `svg.go`. **Nunca** un
  `ssr.go`.
- **Sin `fmt`, `errors`, `strconv`, `strings`, `log`** de la stdlib: usa
  `github.com/tinywasm/fmt`.
- **Sin `encoding/json`, sin `reflect`.**
- **Sin `map[K]V`** en el camino que compila para wasm.

### 0.2 Funciones puras, sin estado

Mismo `Content`, mismos bytes, siempre. Sin variables globales mutables, sin
lectura de reloj, sin aleatoriedad, sin acceso a red ni a disco. La entrada es un
`Content` en memoria y nada más.

### 0.3 Este repositorio no valida ni autoriza

Validar el contenido es trabajo de `veltylabs/site_content` v0.1.0; decidir quién puede
editar es de `veltylabs/site_manager` v0.1.0. Aquí **sólo se dibuja**.

Si un campo llega vacío, se omite la parte que lo usaría — **no** se inventa un
valor por defecto y **no** se entra en pánico. El único pánico legítimo es el que
`layout/landing` ya lanza por metadatos duplicados, y ese se deja pasar a
propósito (§3).

### 0.4 No inventes campos de contenido

Si falta un dato, se agrega en `veltylabs/site_content` y se publica. Este
repositorio **consume** el esquema; no lo extiende ni lo duplica.

### 0.5 Sin strings mágicos, idioma y estructura

Todo string repetido es una constante nombrada. Código en inglés; documentación
y comentarios de prosa en español. Jerarquía plana, archivos de menos de 500
líneas, tests bajo `tests/`.

---

## 1. La API pública — tres funciones

```go
// Landing arma la pagina del sitio a partir de su contenido.
// La composicion es FIJA: el orden de secciones no depende del contenido.
func Landing(c sitecontent.Content, domain string) *landing.Page

// JSONLD produce el documento schema.org del rubro declarado en c.SEO.
func JSONLD(c sitecontent.Content) string

// LLMsTxt produce el resumen en markdown que leen los motores generativos.
func LLMsTxt(c sitecontent.Content) string
```

Nada más se exporta en v1.

> **Desviación respecto a la versión original de este plan:** `Landing` recibe
> `domain` como segundo parámetro. El plan original asumía `c.SEO.Domain`,
> pero ese campo no existe en el `site_content` v0.1.0 publicado. `domain` es
> un dato de infraestructura (qué dominio sirve este sitio), no contenido que
> el cliente edita, así que se pasa explícito en vez de inventarse en
> `site_content` — mantiene `Landing` pura (mismos argumentos, mismos bytes).

---

## 2. `Landing` — `landing.go`

Construye un `*landing.Page` con la composición **fija**, en este orden exacto:

| # | Sección de `landing` | De dónde sale |
|---|---|---|
| 1 | `InfoBar` | `c.Contact` |
| 2 | `Header` | `c.Brand` + menú derivado de las secciones presentes |
| 3 | `Hero` | `c.Hero` |
| 4 | `Split` | `c.About` |
| 5 | `Cards` | `c.Services` |
| 6 | `Stats` | `c.Stats` |
| 7 | `Form` | formulario de contacto |
| 8 | `Hours` | `c.Hours` |
| 9 | `Map` | `c.Map` |
| 10 | `Footer` | `c.Brand` + `c.Contact` |

**El orden no es configurable y no debe volverse configurable.** El cliente edita
contenido, no disposición: un layout que el cliente puede reordenar es un layout
que Velty tiene que reparar gratis.

Una sección cuyo dato está vacío **se omite**, y su ancla desaparece del menú del
`Header`. Un menú que apunta a un ancla que no existe es un enlace roto en cada
página.

### 2.1 Subpáginas

Cada `Service` genera un `landing.SubPage` en `/servicios/<slug>/`, con:

- `Doc.Title` = `Service.Title`
- `Doc.Description` = `Service.Description`
- `Doc.JSONLD` = el bloque schema.org de ese servicio
- `Doc.Canonical` = URL absoluta en el dominio del sitio

**`Title` y `Description` DEBEN ser distintos entre páginas.** `layout/landing`
entra en pánico si se repiten; ese pánico es deseado y no se silencia (§3).

### 2.2 Imágenes

`sitecontent.ImageRef.Key` es una **clave de R2**, no una URL. La ruta final la
construye este paquete, con una constante:

```go
const ImagePathPrefix = "/img/"
```

`Key` → `ImagePathPrefix + Key`. **Nunca** emitas una URL absoluta a R2: el sitio
sirve sus imágenes desde su propio dominio, que es la mitad del argumento de SEO.

---

## 3. El pánico de `landing` se deja pasar

`landing.RenderPages()` entra en pánico si dos páginas comparten `Title` o
`Description`. **No lo captures, no lo conviertas en error, no lo evites
generando sufijos.**

Es la red de seguridad del build de publicación: en CI ese pánico rompe el build
**antes** del commit, y el sitio malo nunca llega al repositorio. Convertirlo en
un `-2` silencioso publicaría dos páginas compitiendo entre sí, que es
exactamente el modo más común de perder ranking.

Quien impide que ocurra es `site_content`, validando al guardar. Aquí sólo se
deja fallar.

---

## 4. `JSONLD` — `jsonld.go`

Produce un documento schema.org según `c.SEO.SchemaType` (`LocalBusiness`,
`MedicalClinic`, …), con nombre, teléfono, dirección, horarios y servicios.

`html.DocumentOptions.JSONLD` es un **string** y se emite verbatim, así que este
paquete es responsable de que sea **JSON válido**. Sin `encoding/json`
disponible, la construcción es manual y el escapado es obligatorio:

```go
// jsonString escribe un string JSON valido: comillas, barras invertidas y
// caracteres de control escapados. Sin esto, una descripcion con comillas
// rompe el JSON-LD entero y el buscador lo descarta en silencio — el peor
// modo de fallo posible, porque el sitio se ve bien.
func jsonString(s string) string
```

**Anti-footgun:** este es el punto exacto donde un escapado a medias produce un
fallo invisible. El caso 6 de §7 existe por eso: una descripción con comillas
dobles, un salto de línea y una barra invertida.

---

## 5. `LLMsTxt` — `llmstxt.go`

Markdown en la raíz del sitio, con: qué es el negocio, dónde está, cómo se le
contacta, qué servicios ofrece y en qué horario.

Vive aquí y no en `tinywasm/sitec` porque su contenido depende del **negocio** —
rubro, servicios, horarios, zona de atención— y `sitec` no sabe qué es un
servicio. Si más adelante todo sitio del ecosistema lo quiere, se promueve a
`sitec` con el contrato ya probado.

Formato estable y determinista: mismo contenido, mismos bytes. Un `llms.txt` que
cambia de orden entre builds produce un commit espurio en cada publicación.

---

## 6. Estructura de archivos

```
landing.go     // Landing(): composicion fija
subpage.go     // Service -> SubPage
jsonld.go      // JSONLD() + jsonString()
llmstxt.go     // LLMsTxt()
theme.go       // constantes compartidas (ImagePathPrefix, rutas, anclas)
docs/ARCHITECTURE.md
tests/
```

Borra `sitetheme.go` (el archivo vacío del andamiaje). Verificación:
`ls sitetheme.go` → no existe.

---

## 7. Tests — `tests/`

| # | Caso | Espera |
|---|---|---|
| 1 | contenido válido completo → `Landing().RenderPages()` | una página por `/` y una por servicio |
| 2 | títulos y descripciones de las páginas | todos distintos |
| 3 | dos servicios con la misma descripción | **pánico** — comportamiento deseado, verifícalo con `recover` |
| 4 | `c.About` vacío | la sección se omite y su ancla no aparece en el menú |
| 5 | `c.Services` vacío | sin subpáginas, sin sección de tarjetas, sin ancla |
| 6 | descripción con `"`, `\` y salto de línea | el JSON-LD sigue siendo JSON válido |
| 7 | `ImageRef.Key` = `foto.webp` | la ruta emitida es `/img/foto.webp`, **sin** dominio |
| 8 | `Landing()` dos veces con el mismo contenido | HTML byte a byte idéntico |
| 9 | `LLMsTxt()` dos veces con el mismo contenido | bytes idénticos |
| 10 | contenido mínimo (sólo campos obligatorios) | no entra en pánico y produce una página |

Los casos 8 y 9 son los que protegen la propiedad del §"La propiedad que hay que
preservar". Escríbelos aunque parezcan triviales: son los que fallan el día que
alguien mete un mapa o una fecha en el camino.

El caso 3 verifica **que sí falle**. Es deliberado.

---

## 8. Documentación

- `docs/ARCHITECTURE.md` — el mapeo contenido → plantilla, la tabla de
  composición fija del §2, y por qué este repositorio vive fuera de
  `veltylabs/modules/`. **Sin código de implementación.**
- `README.md` — inicio rápido, las tres funciones, ejemplo de uso desde los dos
  consumidores. Enlaza todo lo de `docs/` **excepto** `PLAN.md`, que es efímero.
- Si escribes diagramas: **nunca uses `subgraph`** (rompe el renderizado en el
  TUI). `flowchart TD` y `<br/>` para los saltos.

---

## 9. Criterios de aceptación

- [ ] `go vet ./...` limpio; `go test ./tests/...` en verde con los 10 casos.
- [ ] Compila para wasm: `GOOS=js GOARCH=wasm go build ./...` sin errores.
- [ ] `grep -rn "encoding/json\|\"reflect\"\|\"strings\"\|\"errors\"\|\"strconv\"\|\"fmt\"\|\"log\"" --include=*.go . | grep -v _test.go` → vacío.
- [ ] `grep -rn "map\[" --include=*.go . | grep -v _test.go` → vacío.
- [ ] `grep -rn "time.Now\|rand\.\|os.Getenv\|http.Get" --include=*.go .` → vacío: **funciones puras**.
- [ ] `grep -rn "^var [a-z].* = " --include=*.go . | grep -v "_test.go"` → sin estado global mutable.
- [ ] `grep -rn "https://\|r2.dev" --include=*.go . | grep -v _test.go` → vacío: **las imágenes son rutas relativas**.
- [ ] `grep -rn "recover()" --include=*.go . | grep -v _test.go` → vacío: **el pánico de `landing` no se captura**.
- [ ] `grep -rn "internal/" .` → vacío.
- [ ] `ls sitetheme.go` → no existe.
- [ ] `grep -rn "subgraph" docs/` → vacío.

## 10. Fuera de alcance

Una segunda plantilla (se agrega cuando exista un rubro que `landing` no sirva);
validación de contenido (`veltylabs/site_content`); sitios, membresías y planes
(`veltylabs/site_manager`); y cualquier acceso a base de datos, red o disco.
