package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"newsWeb/models"
)

type UserController struct {
	beego.Controller
}

func (this *UserController)ShowRegister()  {
	this.TplName = "register.html"
}

func (this *UserController)HandleRegister()  {
	userName:=this.GetString("userName")
	password:=this.GetString("password")
	//校验数据
	beego.Info(userName,password)
	if userName == ""||password == ""{
		beego.Error("用户名或者密码不能为空")
		this.TplName = "register.html"
		return
	}

	//插入数据
	o:=orm.NewOrm()

	user:=&models.User{
		UserName:userName,
		Pwd:password,
	}

	_,err:=o.Insert(user)
	if err!=nil{
		beego.Error("插入数据失败",err)
		this.TplName = "register.html"
	}
	this.Redirect("/login",302)
}

func (this *UserController)ShowLogin(){
	cookie:=this.Ctx.GetCookie("userName")
	this.Data["cookie"] = cookie
	this.TplName="login.html"
}

func (this *UserController)HandleLogin()  {
	userName:=this.GetString("userName")
	password:=this.GetString("password")
	if userName==""||password==""{
		this.Data["errMsg"]="用户名或者密码不能为空"
		this.TplName = "login.html"
		return
	}
	//get数据
	o:=orm.NewOrm()
	user:=&models.User{
		UserName:userName,
	}
	err:=o.Read(user,"UserName")
	if err!=nil{
		this.Data["errMsg"]="用户名不存在"
		this.TplName = "login.html"
		return
	}
	if user.Pwd!=password{
		this.Data["errMsg"]="账号或密码错误"
		this.TplName = "login.html"
		return
	}
	remember:=this.GetString("remember")
	if remember=="on"{
		this.Ctx.SetCookie("userName",userName,3600*24)
	}else{
		this.Ctx.SetCookie("userName","",-1)
	}

	this.SetSession("userName",userName)
	this.Redirect("/article/index",302)
}

func (this *UserController)Logout()  {
	this.DelSession("userName")
	this.Redirect("/article/index",302)
}