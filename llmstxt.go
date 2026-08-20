package sitetheme

import (
	"github.com/tinywasm/fmt"
	"github.com/veltylabs/site_content"
)

func LLMsTxt(c sitecontent.Content) string {
	b := fmt.GetConv()

	if c.Brand.Name != "" {
		b.WriteString("# ")
		b.WriteString(c.Brand.Name)
		b.WriteString("\n\n")
	}

	if c.Seo.Description != "" {
		b.WriteString("> ")
		b.WriteString(c.Seo.Description)
		b.WriteString("\n\n")
	}

	if c.Contact.Phone != "" || c.Contact.Email != "" || c.Contact.Address != "" {
		b.WriteString("## Contacto\n\n")
		if c.Contact.Address != "" {
			b.WriteString("- Dirección: ")
			b.WriteString(c.Contact.Address)
			b.WriteString("\n")
		}
		if c.Contact.Phone != "" {
			b.WriteString("- Teléfono: ")
			b.WriteString(c.Contact.Phone)
			b.WriteString("\n")
		}
		if c.Contact.Email != "" {
			b.WriteString("- Email: ")
			b.WriteString(c.Contact.Email)
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	if len(c.Services) > 0 {
		b.WriteString("## Servicios\n\n")
		for _, svc := range c.Services {
			b.WriteString("### ")
			b.WriteString(svc.Title)
			b.WriteString("\n")
			if svc.Description != "" {
				b.WriteString(svc.Description)
				b.WriteString("\n")
			}
			if svc.Body != "" {
				b.WriteString(svc.Body)
				b.WriteString("\n")
			}
			b.WriteString("\n")
		}
	}

	if len(c.Hours) > 0 {
		b.WriteString("## Horarios de Atención\n\n")
		for _, sch := range c.Hours {
			b.WriteString("- ")
			b.WriteString(sch.Days)
			b.WriteString(": ")
			b.WriteString(sch.Hours)
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	return b.String()
}
