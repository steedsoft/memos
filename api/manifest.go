package api

type Icon struct {
	Src   string `json:"src"`
	Type  string `json:"type"`
	Sizes string `json:"sizes"`
}

type Manifest struct {
	Name            string `json:"name"`
	ShortName       string `json:"short_name"`
	Description     string `json:"description"`
	Icons           []Icon `json:"icons"`
	StartURL        string `json:"start_url"`
	Scope           string `json:"scope"`
	Display         string `json:"display"`
	ThemeColor      string `json:"theme_color"`
	BackgroundColor string `json:"background_color"`
}
