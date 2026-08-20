package sitetheme

import (
	"github.com/tinywasm/fmt"
	"github.com/veltylabs/site_content"
)

func jsonString(s string) string {
	b := fmt.GetConv()
	b.WriteByte('"')
	fmt.JSONEscape(s, b)
	b.WriteByte('"')
	return b.String()
}

func JSONLD(c sitecontent.Content) string {
	b := fmt.GetConv()

	schemaType := c.Seo.SchemaType
	if schemaType == "" {
		schemaType = "LocalBusiness"
	}

	b.WriteString("{\n")
	b.WriteString("  \"@context\": \"https://schema.org\",\n")
	b.WriteString("  \"@type\": ")
	b.WriteString(jsonString(schemaType))
	b.WriteString(",\n")

	b.WriteString("  \"name\": ")
	b.WriteString(jsonString(c.Brand.Name))

	if c.Contact.Phone != "" {
		b.WriteString(",\n  \"telephone\": ")
		b.WriteString(jsonString(c.Contact.Phone))
	}

	if c.Contact.Email != "" {
		b.WriteString(",\n  \"email\": ")
		b.WriteString(jsonString(c.Contact.Email))
	}

	if c.Seo.Description != "" {
		b.WriteString(",\n  \"description\": ")
		b.WriteString(jsonString(c.Seo.Description))
	}

	if c.Brand.WideLogo != "" {
		b.WriteString(",\n  \"image\": ")
		b.WriteString(jsonString(ImagePathPrefix + c.Brand.WideLogo))
	} else if c.Brand.CompactLogo != "" {
		b.WriteString(",\n  \"image\": ")
		b.WriteString(jsonString(ImagePathPrefix + c.Brand.CompactLogo))
	}

	if c.Contact.Address != "" {
		b.WriteString(",\n  \"address\": {\n")
		b.WriteString("    \"@type\": \"PostalAddress\",\n")
		b.WriteString("    \"streetAddress\": ")
		b.WriteString(jsonString(c.Contact.Address))
		b.WriteString("\n  }")
	}

	if len(c.Hours) > 0 {
		b.WriteString(",\n  \"openingHoursSpecification\": [\n")
		for i, sch := range c.Hours {
			if i > 0 {
				b.WriteString(",\n")
			}
			b.WriteString("    {\n")
			b.WriteString("      \"@type\": \"OpeningHoursSpecification\",\n")
			b.WriteString("      \"dayOfWeek\": ")
			b.WriteString(jsonString(sch.Days))
			b.WriteString(",\n      \"opens\": ")
			b.WriteString(jsonString(sch.Hours))
			b.WriteString("\n    }")
		}
		b.WriteString("\n  ]")
	}

	if len(c.Services) > 0 {
		b.WriteString(",\n  \"hasOfferCatalog\": {\n")
		b.WriteString("    \"@type\": \"OfferCatalog\",\n")
		b.WriteString("    \"name\": \"Servicios\",\n")
		b.WriteString("    \"itemListElement\": [\n")
		for i, svc := range c.Services {
			if i > 0 {
				b.WriteString(",\n")
			}
			b.WriteString("      {\n")
			b.WriteString("        \"@type\": \"Offer\",\n")
			b.WriteString("        \"itemOffered\": {\n")
			b.WriteString("          \"@type\": \"Service\",\n")
			b.WriteString("          \"name\": ")
			b.WriteString(jsonString(svc.Title))
			if svc.Description != "" {
				b.WriteString(",\n          \"description\": ")
				b.WriteString(jsonString(svc.Description))
			}
			b.WriteString("\n        }\n")
			b.WriteString("      }")
		}
		b.WriteString("\n    ]\n")
		b.WriteString("  }")
	}

	b.WriteString("\n}")
	return b.String()
}
