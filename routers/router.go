package routers

import (
	"CDN_Refresh/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
    beego.Router("/refresh", &controllers.MainController{},"post:Refresh")
}
