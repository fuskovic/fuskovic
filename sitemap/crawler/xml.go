package crawler

import (
	"encoding/xml"
	"fmt"
)

const (
	sitemapHeader = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">"
	sitemapFooter = "</urlset>"
)

type (
	url struct {
		Loc string `xml:"loc"`
	}

	urls []url
)

func toURLs(s []string) (URLs urls) {
	for _, link := range s {
		URLs = append(URLs, url{Loc: link})
	}
	return
}

// ToXML turns a slice of strings sitemap protocol compatible xml byte data.
func ToXML(links []string) ([]byte, error) {
	data, err := xml.MarshalIndent(toURLs(links), "\t", "\t")
	if err != nil {
		return []byte{}, err
	}
	xmlData := []byte(fmt.Sprintf("%s\n%s\n%s", sitemapHeader, data, sitemapFooter))
	return xmlData, nil
}
