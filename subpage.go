package sitetheme

import (
	"github.com/tinywasm/fmt"
	"github.com/tinywasm/html"
	"github.com/tinywasm/layout/landing"
	"github.com/veltylabs/site_content"
)

func buildSubPage(svc sitecontent.Service, c sitecontent.Content) landing.SubPage {
	slug := svc.Slug
	if slug == "" {
		slug = fmt.ToLower(svc.Title)
		slug = fmt.ReplaceAll(slug, " ", "-")
	}

	path := fmt.Sprintf("%s%s/", ServicesPathPrefix, slug)

	var canonical string
	if c.SEO.Domain != "" {
		canonical = fmt.Sprintf("https://%s%s", c.SEO.Domain, path)
	}

	doc := html.DocumentOptions{
		Title:       svc.Title,
		Description: svc.Description,
		JSONLD:      JSONLD(c),
		Canonical:   canonical,
	}

	if svc.Image.Key != "" {
		doc.Image = ImagePathPrefix + svc.Image.Key
	}

	imgSrc := ""
	if svc.Image.Key != "" {
		imgSrc = ImagePathPrefix + svc.Image.Key
	}

	var sections []*landing.Section
	splitSec := landing.Split(svc.Title, imgSrc, svc.Description)
	sections = append(sections, splitSec)

	return landing.SubPage{
		Path:     path,
		Doc:      doc,
		Sections: sections,
	}
}
