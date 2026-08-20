# AGENTS.md — `veltylabs/sitetheme`

Guía obligatoria para cualquier agente (o persona) que trabaje en este repo.

## Qué es

El **tema** de los sitios de cliente de Velty: mapea el contenido editable
(`veltylabs/site_content`) a una plantilla concreta
(`github.com/tinywasm/layout/landing`) y produce los artefactos de SEO/GEO que
dependen del negocio — JSON-LD y `llms.txt`.

Lo consumen **dos** lados, y esa es toda su razón de ser:

- el **panel** (`veltylabs/misitio`), para previsualizar en el navegador;
- el **repositorio de sitios** (`veltylabs/clientsites`), para construir en CI.

Como ambos usan el mismo código, **la vista previa no puede desviarse de lo que
se publica**. Esa propiedad es el producto de este repositorio; cualquier cambio
que la rompa es un error, aunque compile.

## Por qué NO vive en `veltylabs/modules/`

Un módulo de dominio de `veltylabs/modules/*` tiene **prohibido** importar un
renderizador concreto: debe poder ensamblarse bajo un renderizador que nunca vio.

Este repositorio hace exactamente lo contrario — **nombrar `layout/landing` es su
trabajo**. Por eso vive fuera de esa carpeta y no obedece esa lista blanca.

El corte importa: `site_content` es el **dato** (agnóstico, entra al Worker);
`sitetheme` es el **tema** (concreto, nunca entra al Worker). Cuando exista una
segunda plantilla se agrega aquí, sin tocar el panel ni la base.

## Restricciones duras

| Regla | Detalle |
|---|---|
| **No entra al Worker** | Nada de este repositorio se importa desde `edge/` de `misitio`. Arrastra el kit de UI entero y el Worker tiene un límite duro de 1 MB. |
| **Compila para wasm** | El panel lo usa en el navegador, así que `landing.go` y equivalentes **no** llevan `//go:build !wasm`. Sólo CSS/SVG/JS/HTML pesado va en archivos `!wasm`. |
| **Sin stdlib pesada** | `tinywasm/fmt` en vez de `fmt`/`errors`/`strings`/`strconv`. Sin `encoding/json`, sin `reflect`. |
| **Sin mapas en el camino wasm** | Slices + búsqueda lineal. Los conjuntos son chicos. |
| **Sin estado** | Las funciones son puras: mismo contenido, mismos bytes. Un tema con estado hace que la vista previa y el build difieran. |
| **Sin `internal/`** | Señal de fork o duplicado de una dependencia en vez de contribuir aguas arriba. |

## No hagas

- **No inventes campos de contenido.** Si falta un dato, se agrega en
  `veltylabs/site_content` y se publica. Este repositorio **consume** el
  esquema, no lo extiende.
- **No metas reglas de negocio.** Validar es trabajo de `site_content`; decidir
  quién puede editar es de `site_manager`. Aquí sólo se dibuja.
- **No leas de la base ni de la red.** La entrada es un `Content` en memoria.
- **No versiones los documentos.** Nada de `v1`/`v2` ni historiales dentro de los
  archivos: eso lo guarda Git.
- **No enlaces `PLAN.md` desde un documento permanente.** Los planes se borran al
  cerrar el ciclo.

## Convención de idioma

Código en inglés (structs, campos, funciones, paquetes). Documentación,
comentarios de prosa y etiquetas de diagramas, en español.

## Estructura y pruebas

- Jerarquía plana; archivos de más de 500 líneas se subdividen por dominio.
- Los productores SSR van en archivos con nombre de extensión y `//go:build !wasm`:
  `css.go`, `js.go`, `html.go`, `svg.go`. **Nunca** en un `ssr.go`.
- Todos los tests bajo `tests/`, contra los paquetes reales.

## Dependencias

| Librería | Rol |
|---|---|
| `github.com/veltylabs/site_content` | el contenido de entrada |
| `github.com/tinywasm/layout/landing` | la plantilla de salida |
| `github.com/tinywasm/html` | `Page`, `DocumentOptions` (incluye el campo `JSONLD`) |
| `github.com/tinywasm/fmt` | reemplazo de `fmt`/`errors`/`strings` |
