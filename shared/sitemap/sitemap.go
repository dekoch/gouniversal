package sitemap

import (
	"fmt"
	"strings"
)

type page struct {
	Menu  string
	Path  string
	Title string
}

type Sitemap struct {
	Pages []page
}

// Register adds a pagepath + title to a sitemap
// e.g.
// pagepath = "App:Program:MyApp"
// title = "MyApp"
func (site *Sitemap) Register(menu string, path string, title string) {

	newpage := make([]page, 1)

	newpage[0].Menu = menu
	newpage[0].Path = path
	newpage[0].Title = title

	site.Pages = append(newpage, site.Pages...)
}

// PageList returns all registered pages
func (site *Sitemap) PageList() []string {

	var list []string
	path := make([]string, 1)

	for i := len(site.Pages) - 1; i >= 0; i-- {

		path[0] = site.Pages[i].Path

		list = append(list, path...)
	}

	return list
}

// PageTitle returns a page title from selected pagepath
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

// ShowMap lists all registered pages
func (site *Sitemap) ShowMap() {

	for i := 0; i < len(site.Pages); i++ {

		fmt.Print(site.Pages[i].Path)
		fmt.Print("\t")
		fmt.Println(site.Pages[i].Title)
	}
}
