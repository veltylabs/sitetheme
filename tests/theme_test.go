package tests

import (
	"encoding/json"
	"testing"

	"github.com/tinywasm/fmt"
	"github.com/veltylabs/site_content"
	"github.com/veltylabs/sitetheme"
)

func sampleContent() sitecontent.Content {
	return sitecontent.Content{
		Brand: sitecontent.Brand{
			Name: "Clínica San José",
			WideLogo: sitecontent.ImageRef{
				Key: "logo-wide.webp",
				Alt: "Logo Ancho",
			},
			CompactLogo: sitecontent.ImageRef{
				Key: "logo-compact.webp",
				Alt: "Logo Compacto",
			},
			LogoAlt: "Clínica San José Logo",
			Href:    "/",
		},
		Contact: sitecontent.Contact{
			Phone:   "+56 9 1234 5678",
			Email:   "contacto@clinicasanjose.cl",
			Address: "Av. Libertad 123, Chillán",
			Hours:   "08:00 - 20:00",
		},
		Hero: sitecontent.Hero{
			Title:    "Cuidamos tu salud y la de tu familia",
			Subtitle: "Atención médica integral con profesionales altamente calificados.",
			CTAs: []sitecontent.Link{
				{Label: "Nuestros Servicios", Href: "#servicios"},
				{Label: "Agendar Cita", Href: "#contacto"},
			},
			Slides: []sitecontent.ImageRef{
				{Key: "hero-1.webp", Alt: "Instalaciones"},
			},
		},
		About: sitecontent.About{
			Title: "Más de 20 años de experiencia",
			Image: sitecontent.ImageRef{
				Key: "nosotros.webp",
				Alt: "Equipo médico",
			},
			Paragraphs: []string{
				"Somos un centro médico comprometido con la excelencia.",
				"Contamos con infraestructura moderna y un equipo dedicado.",
			},
		},
		Services: []sitecontent.Service{
			{
				Title:       "Medicina General",
				Description: "Atención primaria para todas las edades.",
				Slug:        "medicina-general",
				Image:       sitecontent.ImageRef{Key: "medicina-general.webp", Alt: "Medicina General"},
				Badge:       "Popular",
				LinkLabel:   "Ver detalles",
			},
			{
				Title:       "Pediatría",
				Description: "Cuidado especializado para los más pequeños.",
				Slug:        "pediatria",
				Image:       sitecontent.ImageRef{Key: "pediatria.webp", Alt: "Pediatría"},
				Badge:       "Infantil",
				LinkLabel:   "Ver detalles",
			},
		},
		Stats: []sitecontent.Stat{
			{Value: "+10.000", Label: "Pacientes atendidos"},
			{Value: "15", Label: "Especialidades médicas"},
		},
		Hours: sitecontent.Hours{
			Title: "Horarios de Atención",
			Schedules: []sitecontent.Schedule{
				{Days: "Lunes a Viernes", Hours: "08:00 - 20:00"},
				{Days: "Sábados", Hours: "09:00 - 14:00"},
			},
		},
		Map: sitecontent.Map{
			Title: "Nuestra Ubicación",
			URL:   "https://maps.google.com/embed?pb=123",
		},
		SEO: sitecontent.SEO{
			Domain:      "clinicasanjose.cl",
			SchemaType:  "MedicalClinic",
			Title:       "Clínica San José — Atención Médica en Chillán",
			Description: "Centro médico integral en Chillán. Agende su hora en línea.",
		},
	}
}

// Case 1: contenido válido completo -> Landing().RenderPages()
func TestCase1_CompleteValidContent(t *testing.T) {
	c := sampleContent()
	page := sitetheme.Landing(c)
	pages := page.RenderPages()

	// Expected: 1 homepage ("/") and 1 for each service (2) = 3 total pages
	if len(pages) != 3 {
		t.Fatalf("se esperaban 3 páginas, se obtuvieron %d", len(pages))
	}

	if pages[0].Path != "/" {
		t.Errorf("ruta de página principal esperada '/', se obtuvo %s", pages[0].Path)
	}

	if pages[1].Path != "/servicios/medicina-general/" {
		t.Errorf("ruta esperada '/servicios/medicina-general/', se obtuvo %s", pages[1].Path)
	}

	if pages[2].Path != "/servicios/pediatria/" {
		t.Errorf("ruta esperada '/servicios/pediatria/', se obtuvo %s", pages[2].Path)
	}
}

// Case 2: títulos y descripciones de las páginas todos distintos
func TestCase2_DistinctTitlesAndDescriptions(t *testing.T) {
	c := sampleContent()
	page := sitetheme.Landing(c)
	pages := page.RenderPages()

	titles := make([]string, 0, len(pages))
	descs := make([]string, 0, len(pages))

	for _, p := range pages {
		for _, prevTitle := range titles {
			if p.Doc.Title == prevTitle {
				t.Errorf("título duplicado encontrado: %s", p.Doc.Title)
			}
		}
		for _, prevDesc := range descs {
			if p.Doc.Description == prevDesc {
				t.Errorf("descripción duplicada encontrada: %s", p.Doc.Description)
			}
		}
		titles = append(titles, p.Doc.Title)
		descs = append(descs, p.Doc.Description)
	}
}

// Case 3: dos servicios con la misma descripción -> pánico deseado
func TestCase3_DuplicateServiceDescriptionPanics(t *testing.T) {
	c := sampleContent()
	c.Services[1].Description = c.Services[0].Description // Duplicate description

	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("se esperaba pánico por descripciones duplicadas, pero no ocurrió")
		}
	}()

	page := sitetheme.Landing(c)
	_ = page.RenderPages()
}

// Case 4: c.About vacío -> sección omitida y ancla no aparece en el menú
func TestCase4_EmptyAboutOmitted(t *testing.T) {
	c := sampleContent()
	c.About = sitecontent.About{} // Empty About

	page := sitetheme.Landing(c)
	pages := page.RenderPages()

	body := pages[0].Body
	if fmt.Contains(body, `id="nosotros"`) {
		t.Errorf("el ancla 'nosotros' no debería estar presente en el HTML")
	}
	if fmt.Contains(body, "Nosotros") && fmt.Contains(body, `href="#nosotros"`) {
		t.Errorf("el menú no debería contener el enlace a #nosotros")
	}
}

// Case 5: c.Services vacío -> sin subpáginas, sin sección de tarjetas, sin ancla
func TestCase5_EmptyServicesOmitted(t *testing.T) {
	c := sampleContent()
	c.Services = nil // Empty Services

	page := sitetheme.Landing(c)
	pages := page.RenderPages()

	if len(pages) != 1 {
		t.Fatalf("se esperaba solo 1 página (raíz), se obtuvieron %d", len(pages))
	}

	body := pages[0].Body
	if fmt.Contains(body, `id="servicios"`) {
		t.Errorf("el ancla 'servicios' no debería estar presente")
	}
}

// Case 6: descripción con ", \ y salto de línea -> JSON-LD válido
func TestCase6_JSONLDComplexEscaping(t *testing.T) {
	c := sampleContent()
	c.SEO.Description = "Atención con \"comillas\", \\barra invertida\\ y\nsalto de línea."

	jsonldStr := sitetheme.JSONLD(c)

	var parsed map[string]any
	err := json.Unmarshal([]byte(jsonldStr), &parsed)
	if err != nil {
		t.Fatalf("JSON-LD generado no es JSON válido: %v\nString:\n%s", err, jsonldStr)
	}
}

// Case 7: ImageRef.Key = foto.webp -> /img/foto.webp sin dominio
func TestCase7_ImageRelativePath(t *testing.T) {
	c := sampleContent()
	c.About.Image = sitecontent.ImageRef{Key: "foto.webp", Alt: "Foto"}

	page := sitetheme.Landing(c)
	pages := page.RenderPages()

	body := pages[0].Body
	if !fmt.Contains(body, "/img/foto.webp") {
		t.Errorf("se esperaba la ruta de imagen '/img/foto.webp' en el cuerpo")
	}
	if fmt.Contains(body, "https://r2.dev") || fmt.Contains(body, "r2.dev") {
		t.Errorf("no deben emitirse URLs absolutas de R2")
	}
}

// Case 8: Landing() dos veces con el mismo contenido -> HTML idéntico
func TestCase8_LandingPureFunction(t *testing.T) {
	c := sampleContent()

	page1 := sitetheme.Landing(c)
	pages1 := page1.RenderPages()

	page2 := sitetheme.Landing(c)
	pages2 := page2.RenderPages()

	if len(pages1) != len(pages2) {
		t.Fatalf("número de páginas difiere entre ejecuciones")
	}

	for i := range pages1 {
		if pages1[i].Body != pages2[i].Body {
			t.Errorf("página %d difiere entre ejecuciones consecutivas", i)
		}
	}
}

// Case 9: LLMsTxt() dos veces con el mismo contenido -> bytes idénticos
func TestCase9_LLMsTxtPureFunction(t *testing.T) {
	c := sampleContent()

	res1 := sitetheme.LLMsTxt(c)
	res2 := sitetheme.LLMsTxt(c)

	if res1 != res2 {
		t.Errorf("LLMsTxt difiere entre ejecuciones consecutivas:\n1: %s\n2: %s", res1, res2)
	}
}

// Case 10: contenido mínimo -> no entra en pánico y produce una página
func TestCase10_MinimalContent(t *testing.T) {
	c := sitecontent.Content{
		Brand: sitecontent.Brand{
			Name: "Negocio Mínimo",
		},
	}

	page := sitetheme.Landing(c)
	pages := page.RenderPages()

	if len(pages) != 1 {
		t.Fatalf("se esperaba 1 página para contenido mínimo, se obtuvieron %d", len(pages))
	}
}
