package sitetheme

import (
	"github.com/tinywasm/fmt"
	"github.com/tinywasm/html"
	"github.com/tinywasm/layout/landing"
	"github.com/veltylabs/site_content"
)

func buildSubPage(svc sitecontent.Service, c sitecontent.Content, domain string) landing.SubPage {
	slug := svc.Slug
	if slug == "" {
		slug = fmt.ToLower(svc.Title)
		slug = fmt.ReplaceAll(slug, " ", "-")
	}

	path := fmt.Sprintf("%s%s/", ServicesPathPrefix, slug)

	doc := html.DocumentOptions{
		Title:       svc.Title,
		Description: svc.Description,
		JSONLD:      JSONLD(c),
	}
	if domain != "" {
		doc.Canonical = "https://" + domain + path
	}

	if svc.Image != "" {
		doc.Image = ImagePathPrefix + svc.Image
	}

	imgSrc := ""
	if svc.Image != "" {
		imgSrc = ImagePathPrefix + svc.Image
	}

	pText := svc.Description
	if svc.Body != "" {
		pText = svc.Body
	}

	var sections []*landing.Section
	splitSec := landing.Split(svc.Title, imgSrc, pText)
	sections = append(sections, splitSec)

	return landing.SubPage{
		Path:     path,
		Doc:      doc,
		Sections: sections,
	}
}
