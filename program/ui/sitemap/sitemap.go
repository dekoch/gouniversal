package sitemap

import (
	"fmt"
	"strings"
)

type page struct {
	Path  string
	Title string
	Depth int
}

type Sitemap struct {
	Pages []page
}

func (site *Sitemap) Register(path string, title string) {

	newpage := make([]page, 1)

	newpage[0].Path = path
	newpage[0].Title = title
	newpage[0].Depth = strings.Count(path, ":")

	site.Pages = append(newpage, site.Pages...)
}

func (site *Sitemap) PageList() []string {

	var list []string
	path := make([]string, 1)

	for i := len(site.Pages) - 1; i >= 0; i-- {

		path[0] = site.Pages[i].Path

		list = append(list, path...)
	}

	return list
}

func (site *Sitemap) PageTitle(path string) string {

	// find page
	for i := 0; i < len(site.Pages); i++ {

		if path == site.Pages[i].Path {

			return site.Pages[i].Title
		}
	}

	// find page with same prefix (path with parameter)
	for i := 0; i < len(site.Pages); i++ {

		if strings.HasPrefix(path, site.Pages[i].Path) {

			return site.Pages[i].Title
		}
	}

	return ""
}

func (site *Sitemap) ShowMap() {

	for i := 0; i < len(site.Pages); i++ {

		fmt.Println(site.Pages[i].Path)
	}
}
