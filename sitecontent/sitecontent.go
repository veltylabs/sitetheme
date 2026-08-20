package sitecontent

type ImageRef struct {
	Key string
	Alt string
}

type Link struct {
	Label string
	Href  string
}

type Brand struct {
	Name        string
	WideLogo    ImageRef
	CompactLogo ImageRef
	LogoAlt     string
	Href        string
}

type Contact struct {
	Phone   string
	Email   string
	Address string
	Hours   string
}

type Hero struct {
	Title    string
	Subtitle string
	CTAs     []Link
	Slides   []ImageRef
}

type About struct {
	Title      string
	Image      ImageRef
	Paragraphs []string
}

type Service struct {
	Title       string
	Description string
	Slug        string
	Image       ImageRef
	Badge       string
	LinkLabel   string
}

type Stat struct {
	Value string
	Label string
}

type Schedule struct {
	Days  string
	Hours string
}

type Hours struct {
	Title     string
	Schedules []Schedule
}

type Map struct {
	Title string
	URL   string
}

type SEO struct {
	Domain     string
	SchemaType string
	Title      string
	Description string
}

type Content struct {
	Brand    Brand
	Contact  Contact
	Hero     Hero
	About    About
	Services []Service
	Stats    []Stat
	Hours    Hours
	Map      Map
	SEO      SEO
}
