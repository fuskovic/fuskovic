package crawler

import "github.com/fuskovic/gophercises/linkparser/linx"

const rootLvl = 1

type (
	linkPath struct {
		src, dest, linkText string
		lvl                 int
	}

	linkPaths []linkPath

	sitemap map[int]linkPaths

	// Crawler maintains the relationships between different levels and their linkPaths.
	Crawler struct {
		baseURL string
		depth   int
		siteMap sitemap
	}
)

// Init initializes a new Crawler with the baseURL for a particular domain and a depth specifying how many levels to crawl.
func Init(baseURL string, depth int) (*Crawler, error) {
	allBasePageLinks, err := getAllLinksFromPage(baseURL, baseURL)
	if err != nil {
		return nil, err
	}
	return &Crawler{
		baseURL: baseURL,
		depth:   depth,
		siteMap: map[int]linkPaths{rootLvl: linksToLinkPaths(rootLvl, baseURL, *allBasePageLinks)},
	}, nil
}

// CrawlByDepth crawls a site going as deep as the depth specified.
func (c *Crawler) CrawlByDepth() (results []string) {
	for i := 1; i < c.depth; i++ {
		for _, lp := range c.siteMap[i] {
			links, err := getAllLinksFromPage(c.baseURL, lp.dest)
			if err != nil {
				continue
			}
			lps := linksToLinkPaths(c.depth, lp.src, filterInvalidLinks(c.baseURL, *links))
			nextLv := i + 1
			c.siteMap[nextLv] = lps
		}
	}

	counts := make(map[string]int)

	for lvl := 1; lvl <= c.depth; lvl++ {
		for _, lp := range c.siteMap[lvl] {
			counts[lp.dest]++

			if counts[lp.dest] == 1 {
				results = append(results, lp.dest)
			}
		}
	}
	return results
}

// CrawlAll crawls a site (regardless of depth) until all links within the domain have been found.
func (c Crawler) CrawlAll() (results []string) {
	var visited, notVisited linx.Links

	for _, lps := range c.siteMap {
		for _, lp := range lps {
			notVisited = append(notVisited, linx.Link{
				Href: lp.dest,
				Text: lp.linkText,
			})
		}
	}

	planningToVisitAlready := func(link linx.Link) bool {
		for _, notVisitedLink := range notVisited {
			if link.Href == notVisitedLink.Href {
				return true
			}
		}
		return false
	}

	markAsVisited := func(link linx.Link) {
		var updated linx.Links
		for _, notVisitedLink := range notVisited {
			if link.Href != notVisitedLink.Href {
				updated = append(updated, link)
			}
		}
		notVisited = updated
		visited = append(visited, link)
	}

	hasBeenVisitedAlready := func(link linx.Link) bool {
		for _, visitedLink := range visited {
			if link.Href == visitedLink.Href {
				return true
			}
		}
		return false
	}

	needsToBeVisited := func(link linx.Link) bool {
		return !planningToVisitAlready(link) && !hasBeenVisitedAlready(link)
	}

	for len(notVisited) > 0 {
		for _, link := range notVisited {
			links, err := getAllLinksFromPage(c.baseURL, link.Href)
			if err != nil {
				continue
			} else {
				markAsVisited(link)
			}

			for _, link := range filterInvalidLinks(c.baseURL, *links) {
				if needsToBeVisited(link) {
					notVisited = append(notVisited, link)
				}
			}
		}
	}

	for _, link := range dedupeLinks(visited) {
		results = append(results, link.Href)
	}
	return
}
