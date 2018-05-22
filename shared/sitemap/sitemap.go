package sitemap

import (
	"fmt"
	"strings"
)

type Page struct {
	Menu       string
	MenuOrder  int
	Path       string
	Title      string
	TitleOrder int
}

type Sitemap struct {
	Pages []Page
}

// RegisterWithOrder adds a pagepath + title to a sitemap
// e.g.
// menu = "Program"
// menuorder = 10
// pagepath = "App:Program:MyApp"
// title = "MyApp"
// titleorder = 10
func (site *Sitemap) RegisterWithOrder(menu string, menuorder int, path string, title string, titleorder int) {

	newPage := make([]Page, 1)

	newPage[0].Menu = menu
	newPage[0].MenuOrder = menuorder
	newPage[0].Path = path
	newPage[0].Title = title
	newPage[0].TitleOrder = titleorder

	site.Pages = append(site.Pages, newPage...)
}

// Register adds a pagepath + title to a sitemap
// e.g.
// menu = "Program"
// pagepath = "App:Program:MyApp"
// title = "MyApp"
func (site *Sitemap) Register(menu string, path string, title string) {

	site.RegisterWithOrder(menu, -1, path, title, -1)
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
