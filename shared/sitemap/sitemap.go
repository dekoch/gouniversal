package sitemap

import (
	"fmt"
	"strings"

	"github.com/dekoch/gouniversal/shared/functions"
)

type Page struct {
	Menu       string
	MenuOrder  int
	Path       string
	Title      string
	TitleOrder int
}

type Sitemap struct {
	pages []Page
}

// RegisterWithOrder adds a pagepath + title to a sitemap
// e.g.
// menu = "Program"
// menuorder = 10
// pagepath = "App:Program:MyApp"
// title = "MyApp"
// titleorder = 10
func (site *Sitemap) RegisterWithOrder(menu string, menuorder int, path string, title string, titleorder int) {

	if functions.IsEmpty(path) ||
		functions.IsEmpty(title) {
		return
	}

	var n Page

	n.Menu = menu
	n.MenuOrder = menuorder
	n.Path = path
	n.Title = title
	n.TitleOrder = titleorder

	site.pages = append(site.pages, n)
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

	for i := len(site.pages) - 1; i >= 0; i-- {

		list = append(list, site.pages[i].Path)
	}

	return list
}

// PageTitle returns a page title from selected pagepath
func (site *Sitemap) PageTitle(path string) string {

	// find page
	for i := 0; i < len(site.pages); i++ {

		if path == site.pages[i].Path {

			return site.pages[i].Title
		}
	}

	// find page with same prefix (path with parameter)
	for i := 0; i < len(site.pages); i++ {

		if strings.HasPrefix(path, site.pages[i].Path) {

			return site.pages[i].Title
		}
	}

	return ""
}

// GetPages returns a array of registered pages
func (site *Sitemap) GetPages() []Page {

	return site.pages
}

// ShowMap lists all registered pages
func (site *Sitemap) ShowMap() {

	for i := 0; i < len(site.pages); i++ {

		fmt.Print(site.pages[i].Path)
		fmt.Print("\t")
		fmt.Println(site.pages[i].Title)
	}
}

// Clear removes all registered pages
func (site *Sitemap) Clear() {

	site.pages = make([]Page, 0)
}
