# Arquitectura de `veltylabs/sitetheme`

## Descripción General

`sitetheme` es el tema estándar de los sitios de clientes en la plataforma Velty. Su función principal es transformar la estructura de contenido editable (`veltylabs/site_content`) en la plantilla `tinywasm/layout/landing`, generando además los artefactos SEO/GEO necesarios (`JSON-LD` y `llms.txt`).

## Consumidores

El código de este repositorio es consumido por dos entornos distintos:

1. **Panel del Cliente (`veltylabs/misitio`)**: Se ejecuta en el navegador (WASM) para previsualización en tiempo real.
2. **Generador de Sitios (`veltylabs/clientsites`)**: Se ejecuta en CI para compilar y publicar las páginas estáticas del sitio.

Ambos consumidores comparten exactamente el mismo paquete para garantizar que la vista previa no se desvíe del sitio publicado final.

## Razones de Ubicación fuera de `veltylabs/modules/`

Los módulos de dominio en `veltylabs/modules/` tienen prohibido importar renderizadores concretos para mantener agnosticismo visual. `sitetheme`, en cambio, tiene como responsabilidad explícita la integración directa con la plantilla `tinywasm/layout/landing`.

- `site_content` = Dato agnóstico
- `sitetheme` = Tema y renderizador concreto

## Mapeo de Contenido a Plantilla (Composición Fija)

El orden de secciones en la página principal es **estricto y fijo**:

```flowchart TD
    InfoBar["1. InfoBar (c.Contact)"]
    Header["2. Header (c.Brand + Menú)"]
    Hero["3. Hero (c.Hero)"]
    Split["4. Split (c.About)"]
    Cards["5. Cards (c.Services)"]
    Stats["6. Stats (c.Stats)"]
    Form["7. Form (Contacto)"]
    Hours["8. Hours (c.Hours)"]
    Map["9. Map (c.Map)"]
    Footer["10. Footer (c.Brand + c.Contact)"]

    InfoBar --> Header
    Header --> Hero
    Hero --> Split
    Split --> Cards
    Cards --> Stats
    Stats --> Form
    Form --> Hours
    Hours --> Map
    Map --> Footer
```

| # | Sección | Fuente de Datos (`sitecontent.Content`) |
|---|---|---|
| 1 | `InfoBar` | `c.Contact` |
| 2 | `Header` | `c.Brand` + menú dinámico según secciones presentes |
| 3 | `Hero` | `c.Hero` |
| 4 | `Split` | `c.About` |
| 5 | `Cards` | `c.Services` |
| 6 | `Stats` | `c.Stats` |
| 7 | `Form` | Formulario de contacto |
| 8 | `Hours` | `c.Hours` |
| 9 | `Map` | `c.Map` |
| 10 | `Footer` | `c.Brand` + `c.Contact` |

Si una sección carece de datos en `Content`, la sección se omite y su ancla de navegación se remueve automáticamente del menú.

## Subpáginas

Cada servicio declarado en `c.Services` genera una subpágina individual bajo la ruta `/servicios/<slug>/`. Cada subpágina cuenta con su título, descripción, canonical URL y bloque schema.org JSON-LD específico.

## El parámetro `domain`

`Landing(c sitecontent.Content, domain string)` recibe el dominio del sitio como segundo argumento, no como parte de `Content`. `site_content` v0.1.0 no tiene un campo de dominio: es un dato de infraestructura (qué sitio sirve este contenido), no contenido editable por el cliente, así que cada consumidor (`misitio`, `clientsites`) lo provee según su propia configuración. Con `domain == ""` no se emite `<link rel="canonical">` ni `og:url` — la función sigue siendo pura y no falla por su ausencia.
