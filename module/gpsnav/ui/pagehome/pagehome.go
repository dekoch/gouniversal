package pagehome

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/dekoch/gouniversal/module/gpsnav/global"
	"github.com/dekoch/gouniversal/module/gpsnav/lang"
	"github.com/dekoch/gouniversal/module/gpsnav/typenav"
	"github.com/dekoch/gouniversal/shared/functions"
	"github.com/dekoch/gouniversal/shared/navigation"
)

func LoadConfig() {
}

func RegisterPage(page *typenav.Page, nav *navigation.Navigation) {

	nav.Sitemap.Register(page.Lang.Home.Menu, "App:GPSNav:Home", page.Lang.Home.Title)
}

func Render(page *typenav.Page, nav *navigation.Navigation, r *http.Request) {

	type content struct {
		Lang    lang.Home
		Bearing template.HTML
		Path    template.HTML
	}
	var c content

	c.Lang = page.Lang.Home

	bearing, _ := global.Geo.GetTargetBearing()

	c.Bearing = template.HTML(strconv.FormatFloat(bearing, 'f', 1, 64))
	c.Path = path(bearing)

	cont, err := functions.PageToString(global.Config.UIFileRoot+"home.html", c)
	if err == nil {
		page.Content += cont
	} else {
		nav.RedirectPath("404", true)
	}
}

func path(bearing float64) template.HTML {

	intBearing := int(functions.Round(bearing, .5, 0))

	ret := `
	<script>
		var canvas = document.getElementById('myCanvas');
		var context = canvas.getContext('2d');
		
		// translate context to center of canvas
		context.translate(canvas.width / 2, canvas.height / 2);

		// arrow
		context.beginPath();
		context.moveTo(0, -50);
		context.lineTo(5, -30);
		context.lineTo(-5, -30);
		context.lineTo(0, -50);
		context.strokeStyle = 'blue';
		context.stroke();
  
		// current direction
		context.beginPath();
		context.moveTo(0, -30);
		context.lineTo(0, 50);
		context.strokeStyle = 'blue';
		context.stroke();
		
		// target direction
		context.rotate(` + strconv.Itoa(intBearing) + ` * Math.PI / 180);
		context.beginPath();
		context.moveTo(0, -50);
		context.lineTo(0, 0);
		context.strokeStyle = 'green';
		context.stroke();
  	</script>`

	return template.HTML(ret)
}
