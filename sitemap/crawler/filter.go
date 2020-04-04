package crawler

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/fuskovic/gophercises/linkparser/linx"
)

func getAllLinksFromPage(baseURL, page string) (*linx.Links, error) {
	resp, err := http.Get(page)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	links, err := linx.GetLinks(data)
	if err != nil {
		return nil, err
	}
	links = filterInvalidLinks(baseURL, links)
	return &links, nil
}

func filterInvalidLinks(baseURL string, links linx.Links) linx.Links {
	var filtered linx.Links

	for _, link := range links {
		if link.Href == "" {
			continue
		}

		if isPath(link.Href) {
			if string(link.Href[0]) == "#" {
				link.Href = "/" + link.Href
			}
			link.Href = baseURL + link.Href
		}

		if isWithinDomain(baseURL, link.Href) {
			filtered = append(filtered, link)
		}
	}
	return dedupeLinks(filtered)
}

func linksToLinkPaths(lvl int, src string, links linx.Links) (lps linkPaths) {
	for _, link := range links {
		lps = append(lps, linkPath{
			src:      src,
			dest:     link.Href,
			linkText: link.Text,
			lvl:      lvl,
		})
	}
	return
}

func dedupeLinks(links linx.Links) (deduped linx.Links) {
	counts := make(map[string]int)
	for _, link := range links {
		counts[link.Href]++

		if counts[link.Href] == 1 {
			deduped = append(deduped, link)
		}
	}
	return
}

func isPath(s string) bool {
	return string(s[0]) == "/" || string(s[0]) == "#"
}

func isWithinDomain(baseURL, s string) bool {
	return strings.Contains(s, baseURL)
}
