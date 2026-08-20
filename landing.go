package sitetheme

import (
	"github.com/tinywasm/fmt"
	"github.com/tinywasm/html"
	"github.com/tinywasm/layout/landing"
	"github.com/veltylabs/site_content"
)

func Landing(c sitecontent.Content) *landing.Page {
	// 1. Convert Brand
	brand := landing.Brand{
		Name:    c.Brand.Name,
		LogoAlt: c.Brand.LogoAlt,
		Href:    c.Brand.Href,
	}
	if c.Brand.WideLogo.Key != "" {
		brand.WideLogoSrc = ImagePathPrefix + c.Brand.WideLogo.Key
	}
	if c.Brand.CompactLogo.Key != "" {
		brand.CompactLogoSrc = ImagePathPrefix + c.Brand.CompactLogo.Key
	}

	// 2. Convert Contact
	contact := landing.Contact{
		Phone:   c.Contact.Phone,
		Email:   c.Contact.Email,
		Address: c.Contact.Address,
		Hours:   c.Contact.Hours,
	}

	// Build menu based on present sections
	var menu []landing.Link

	if c.Hero.Title != "" || c.Hero.Subtitle != "" || len(c.Hero.CTAs) > 0 || len(c.Hero.Slides) > 0 {
		menu = append(menu, landing.Link{Label: NavLabelHero, Href: "#" + AnchorHero})
	}
	if c.About.Title != "" || len(c.About.Paragraphs) > 0 || c.About.Image.Key != "" {
		menu = append(menu, landing.Link{Label: NavLabelAbout, Href: "#" + AnchorAbout})
	}
	if len(c.Services) > 0 {
		menu = append(menu, landing.Link{Label: NavLabelServices, Href: "#" + AnchorServices})
	}
	if len(c.Stats) > 0 {
		menu = append(menu, landing.Link{Label: NavLabelStats, Href: "#" + AnchorStats})
	}
	menu = append(menu, landing.Link{Label: NavLabelForm, Href: "#" + AnchorForm})
	if c.Hours.Title != "" || len(c.Hours.Schedules) > 0 {
		menu = append(menu, landing.Link{Label: NavLabelHours, Href: "#" + AnchorHours})
	}
	if c.Map.Title != "" || c.Map.URL != "" {
		menu = append(menu, landing.Link{Label: NavLabelMap, Href: "#" + AnchorMap})
	}

	// 3. Assemble sections in fixed order
	var sections []*landing.Section

	// 1. InfoBar
	if c.Contact.Phone != "" || c.Contact.Email != "" || c.Contact.Address != "" || c.Contact.Hours != "" {
		sections = append(sections, landing.InfoBar(contact))
	}

	// 2. Header
	sections = append(sections, landing.Header(menu...))

	// 3. Hero
	if c.Hero.Title != "" || c.Hero.Subtitle != "" || len(c.Hero.CTAs) > 0 || len(c.Hero.Slides) > 0 {
		var ctas []landing.Link
		for _, cta := range c.Hero.CTAs {
			ctas = append(ctas, landing.Link{
				Label: cta.Label,
				Href:  cta.Href,
			})
		}
		var slides []landing.Slide
		for _, s := range c.Hero.Slides {
			if s.Key != "" {
				slides = append(slides, landing.Slide{
					Image: ImagePathPrefix + s.Key,
				})
			}
		}
		heroSec := landing.Hero(c.Hero.Title, c.Hero.Subtitle, ctas, slides...).At(AnchorHero)
		sections = append(sections, heroSec)
	}

	// 4. Split (About)
	if c.About.Title != "" || len(c.About.Paragraphs) > 0 || c.About.Image.Key != "" {
		imgSrc := ""
		if c.About.Image.Key != "" {
			imgSrc = ImagePathPrefix + c.About.Image.Key
		}
		splitSec := landing.Split(c.About.Title, imgSrc, c.About.Paragraphs...).At(AnchorAbout)
		sections = append(sections, splitSec)
	}

	// 5. Cards (Services)
	if len(c.Services) > 0 {
		var cards []landing.Card
		for _, svc := range c.Services {
			slug := svc.Slug
			if slug == "" {
				slug = fmt.ToLower(svc.Title)
				slug = fmt.ReplaceAll(slug, " ", "-")
			}
			href := fmt.Sprintf("%s%s/", ServicesPathPrefix, slug)
			imgSrc := ""
			if svc.Image.Key != "" {
				imgSrc = ImagePathPrefix + svc.Image.Key
			}
			cards = append(cards, landing.Card{
				Title:       svc.Title,
				Description: svc.Description,
				Image:       imgSrc,
				Href:        href,
				Badge:       svc.Badge,
				LinkLabel:   svc.LinkLabel,
			})
		}
		cardsSec := landing.Cards(NavLabelServices, cards...).At(AnchorServices)
		sections = append(sections, cardsSec)
	}

	// 6. Stats
	if len(c.Stats) > 0 {
		var stats []landing.Stat
		for _, st := range c.Stats {
			stats = append(stats, landing.Stat{
				Value: st.Value,
				Label: st.Label,
			})
		}
		statsSec := landing.Stats(stats...).At(AnchorStats)
		sections = append(sections, statsSec)
	}

	// 7. Form (always present)
	formSec := landing.Form(FormTitleDefault, FormIntroDefault, nil).At(AnchorForm)
	sections = append(sections, formSec)

	// 8. Hours
	if c.Hours.Title != "" || len(c.Hours.Schedules) > 0 {
		var schedules []landing.Schedule
		for _, sch := range c.Hours.Schedules {
			schedules = append(schedules, landing.Schedule{
				Days:  sch.Days,
				Hours: sch.Hours,
			})
		}
		hoursTitle := c.Hours.Title
		if hoursTitle == "" {
			hoursTitle = NavLabelHours
		}
		hoursSec := landing.Hours(hoursTitle, contact, schedules...).At(AnchorHours)
		sections = append(sections, hoursSec)
	}

	// 9. Map
	if c.Map.Title != "" || c.Map.URL != "" {
		mapTitle := c.Map.Title
		if mapTitle == "" {
			mapTitle = NavLabelMap
		}
		mapSec := landing.MapEmbed(mapTitle, c.Map.URL).At(AnchorMap)
		sections = append(sections, mapSec)
	}

	// 10. Footer
	footerSec := landing.Footer(menu...)
	sections = append(sections, footerSec)

	page := landing.New(brand, sections...)

	// SEO metadata for root page
	var canonical string
	if c.SEO.Domain != "" {
		canonical = fmt.Sprintf("https://%s/", c.SEO.Domain)
	}

	docOpts := html.DocumentOptions{
		Title:       c.SEO.Title,
		Description: c.SEO.Description,
		JSONLD:      JSONLD(c),
		Canonical:   canonical,
	}

	page.WithSEO(docOpts)

	// Add subpages for services
	for _, svc := range c.Services {
		subpage := buildSubPage(svc, c)
		page.AddSubPage(subpage)
	}

	return page
}
