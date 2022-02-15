package packets

type packet struct {
	Name        string
	Url         string
	UrlType     string `yaml:"urlType"`
	Description string
}
