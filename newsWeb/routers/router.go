package routers

import (
	"newsWeb/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
    beego.Router("register",&controllers.UserController{},"get:ShowRegister;post:HandleRegister")
    beego.Router("/login",&controllers.UserController{},"get:ShowLogin;post:HandleLogin")
    beego.Router("/article/index",&controllers.ArticleController{},"get:ShowIndex")
    beego.Router("/article/logout",&controllers.UserController{},"get:Logout")
	beego.Router("/article/add",&controllers.ArticleController{},"get:ShowAdd;post:HandleAdd")
	beego.Router("/article/content",&controllers.ArticleController{},"get:ShowContent")
	beego.Router("/article/update",&controllers.ArticleController{},"get:ShowUpdate;post:HandleUpdate")
    beego.Router("/article/delete",& controllers.ArticleController{},"get:ShowDelete")
	beego.Router("/article/addType",& controllers.ArticleController{},"get:ShowAddType;post:HandleAddType")
    beego.Router("/article/delType",&controllers.ArticleController{},"get:DeleteType")
    beego.Router("/goredis",&controllers.GoRedis{},"get:ShowGet")
}
