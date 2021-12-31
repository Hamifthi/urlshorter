package urlshortner

type Url struct {
	KEY string
	URL string
}

type UrlService interface {
	GetUrl(key string) (*Url, error)
	CreateShortLink(url string) (string, error)
	DeleteShortLink(key string) (string, error)
	DeleteAllMatchingUrls(proposedUrl string) (string, error)
}

type RandomService interface {
	GenerateRandomString() string
}
