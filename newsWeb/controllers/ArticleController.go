package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"math"
	"newsWeb/models"
	"path"
	"time"
)

type ArticleController struct {
	beego.Controller
}

//展示首页
func (this*ArticleController)ShowIndex(){

	userName:=this.GetSession("userName")
	if userName == nil{
		this.Redirect("/login",302)
		return
	}
	//获取数据
	//select * from article;
	//创建orm对象
	o := orm.NewOrm()
	//指定表
	qs :=o.QueryTable("Article") //queryseter  查询对象
	//定义一个容器
	var articles []models.Article

	//获取文章的总数量

	//获取每页的条数
	page := 2
	pageIndex,err:=this.GetInt("pageIndex")
	if err!=nil{
		pageIndex=1
	}
	start:=(pageIndex-1)*2
	var count int64
	selecter:=this.GetString("select")
	if selecter==""||selecter=="全部"{
		count,err=qs.RelatedSel("ArticleType").Count()
		if err!=nil{
			beego.Error(err)
			this.Redirect("/article/index",302)
			return
		}
		_,err=qs.Limit(page,start).RelatedSel("ArticleType").All(&articles)
	}else{
		//orm中一对多的查询是惰性查询
		count,err=qs.RelatedSel("ArticleType").Filter("ArticleType__TypeName",selecter).Count()
		if err!=nil{
			beego.Error(err)
			this.Redirect("/article/index",302)
			return
		}
		_,err=qs.Limit(page,start).RelatedSel("ArticleType").Filter("ArticleType__TypeName",selecter).All(&articles)
	}

	if err!=nil{
		beego.Error(err)
		this.Redirect("/article/index",302)
		return
	}
	var articleType []models.ArticleType
	_,err=o.QueryTable("ArticleType").All(&articleType)
	if err!=nil{
		beego.Error(err)
		this.Redirect("/article/index",302)
		return
	}
	pageCount := int(math.Ceil(float64(count)/float64(page)))
	this.Data["selecter"] = selecter
	this.Data["ArticleType"] =articleType
	this.Data["count"] = count
	this.Data["pageCount"] = pageCount
	//传递数据
	this.Data["articles"] = articles
	this.Data["pageIndex"] = pageIndex
	this.Layout="layout.html"
	this.TplName = "index.html"
}

//展示添加文章页面
func(this*ArticleController)ShowAdd(){
	o:=orm.NewOrm()

	var articleType []models.ArticleType

	_,err:=o.QueryTable("ArticleType").All(&articleType)
	if err!=nil{
		beego.Error(err)
		this.Redirect("/article/index",302)
	}

	this.Data["articleType"] = articleType
	this.Layout="layout.html"
	this.TplName = "add.html"
}

//处理添加文章业务
func(this*ArticleController)HandleAdd(){
	//获取数据
	title :=this.GetString("articleName")
	content :=this.GetString("content")
	file,head,err :=this.GetFile("uploadname")
	articleTypeName:=this.GetString("select")
	defer file.Close()

	//校验数据
	if title == "" || content == "" || err != nil{
		this.Data["errmsg"] = "添加文章失败，请重新添加！"
		this.TplName = "add.html"
		return
	}

	//beego.Info(file,head)

	//1.文件存在覆盖的问题
	//加密算法
	//当前时间
	fileName := time.Now().Format("2006-01-02-15-04-05")
	ext := path.Ext(head.Filename)
	beego.Info(head.Filename,ext)
	//2.文件类型也需要校验
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg"{
		this.Data["errmsg"] = "上传图片格式不正确，请重新上传"
		this.TplName = "add.html"
		return
	}
	//3.文件大小校验
	if head.Size > 5000000 {
		this.Data["errmsg"] = "上传图片过大，请重新上传"
		this.TplName = "add.html"
		return
	}

	//把图片存起来
	this.SaveToFile("uploadname","./static/img/"+fileName+ext)

	//处理数据
	//数据库的插入操作
	//获取orm对象
	o := orm.NewOrm()
	//获取一个插入对象
	var articleType models.ArticleType
	articleType.TypeName = articleTypeName

	err=o.Read(&articleType,"TypeName")
	if err!=nil{
		beego.Error(err)
		this.Redirect("/article/index",302)
		return
	}

	var article models.Article
	//给插入对象赋值
	article.Title = title
	article.Content = content
	article.ArticleType = &articleType
	article.Img = "/static/img/"+fileName+ext
	//插入到数据库
	o.Insert(&article)

	//返回数据
	this.Redirect("/article/index",302)
}


//查看文章详情
func (this *ArticleController)ShowContent()  {
	articleId,err:=this.GetInt("articleId")
	if err!=nil{
		beego.Error("请求连接错误：",err)
		this.Redirect("/article/index",302)
		return
	}

	o:=orm.NewOrm()
	var article models.Article
	article.Id = articleId
	err=o.Read(&article)
	if err!=nil{
		beego.Error("查询文章不存在")
		this.Redirect("/article/index",302)
		return
	}
	var user models.User
	userName:=this.GetSession("userName")
	if userName==nil{
		this.Redirect("/login",302)
		return
	}
	user.UserName = userName.(string)
	o.Read(&user,"UserName")
	//多对多的添加
	m2m:=o.QueryM2M(&article,"Users")
	m2m.Add(user)
	//多对多查询
	//_,err=o.LoadRelated(&article,"Users")
	o.QueryTable("User").RelatedSel()
	qs:=o.QueryTable("User")
	var users []models.User
	qs.Filter("Articles__Article__Id",article.Id).Distinct().All(&users)
	//qs.Filter("Articles__Article__Id",article.Id).Distinct().All(&users)
	this.Data["users"] = users
	this.Data["article"] = article
	this.Layout="layout.html"
	this.TplName = "content.html"

}

//展示编辑页面
func (this *ArticleController)ShowUpdate()  {
	articleId,err:=this.GetInt("articleId")
	if err!=nil{
		beego.Error("获取文章信息失败",err,articleId)
		this.Redirect("/article/index",302)
		return
	}
	//更新
	o:= orm.NewOrm()

	article :=&models.Article{
		Id:articleId,
	}
	err=o.Read(article)
	if err!=nil{
		beego.Error("获取文章信息失败",err)
		this.Redirect("/article/index",302)
		return
	}
	this.Data["article"] = article
	this.Layout="layout.html"
	this.TplName = "update.html"
}

func UpLoad(this *ArticleController,filePath string)(string,error)  {
	file,head,err :=this.GetFile(filePath)
	//校验数据
	if err != nil{
		return "",err
	}
	defer file.Close()
	//1.文件存在覆盖的问题
	//加密算法

	//当前时间
	fileName := time.Now().Format("2006-01-02-15-04-05")
	ext := path.Ext(head.Filename)
	beego.Info(head.Filename,ext)
	//2.文件类型也需要校验
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg"{
		return "",err
	}
	//3.文件大小校验
	if head.Size > 5000000 {
		return "",err
	}

	//把图片存起来
	err=this.SaveToFile(filePath,"./static/img/"+fileName+ext)
	if err!=nil{
		return "",err
	}
	return "/static/img/"+fileName+ext,nil
}

//更新文章
func (this *ArticleController)HandleUpdate()  {
	//获取数据信息
	articleName:=this.GetString("articleName")
	content:=this.GetString("content")
	filePath,err1:=UpLoad(this,"uploadname")
	articleId,err2 := this.GetInt("articleId")
	if articleName==""||content==""||err1!=nil||err2!=nil{
		fmt.Println(articleName,content,filePath)
		beego.Error("修改文章失败，请重新添加！",err1,err2)
		this.Data["errmsg"] = "修改文章失败，请重新添加！"
		this.Redirect("/article/index",302)
		return
	}
	//新建orm对象
	o:=orm.NewOrm()
	//新建article对象
	article:=&models.Article{
		Id:articleId,
	}
	err:=o.Read(article)
	if err!=nil{
		beego.Error("修改文章失败，请重新添加！",err)
		this.Data["errmsg"] = "修改文章失败，请重新添加！"
		this.Redirect("/article/index",302)
		return
	}
	article.Title = articleName
	article.Content = content
	if filePath!=""{
		article.Img = filePath
	}
	_,err=o.Update(article)
	if err!=nil{
		beego.Error("修改文章失败，请重新添加！",err)
		this.Data["errmsg"] = "修改文章失败，请重新添加！"
		this.Redirect("/article/index",302)
		return
	}
	this.Redirect("/article/index",302)
}

//删除文章

func (this *ArticleController)ShowDelete()  {
	articleId,err:=this.GetInt("articleId")
	if err!=nil{
		beego.Error("删除文章失败",err)
		this.Redirect("/article/index",302)
		return
	}
	o:=orm.NewOrm()
	article:=&models.Article{
		Id: articleId,
	}
	_,err=o.Delete(article)
	if err!=nil{
		beego.Error("删除文章失败")
		this.Redirect("/article/index",302)
		return
	}
	this.Redirect("/article/index",302)
}

/*******************
分类
********************/
func (this *ArticleController)ShowAddType()  {
	o:=orm.NewOrm()
	var articleTypes []models.ArticleType
	if len(articleTypes)==0{
		o.QueryTable("ArticleType").All(&articleTypes)
	}
	this.Data["articleTypes"] = articleTypes
	this.Layout="layout.html"
	this.TplName = "addType.html"
}

func (this *ArticleController)HandleAddType()  {
	typeName:=this.GetString("typeName")

	if typeName == ""{
		beego.Error("文章类型不能为空")
		this.Redirect("/article/addType",302)
		return
	}

	o:=orm.NewOrm()
	var articleType models.ArticleType
	articleType.TypeName = typeName
	_,err:=o.Insert(&articleType)
	if err!=nil{
		beego.Error("文章类型已存在",err)
		this.Redirect("/article/addType",302)
		return
	}
	this.Redirect("/article/addType",302)
}

func (this *ArticleController)DeleteType()  {
	id,err:=this.GetInt("typeId")
	if err!=nil{
		beego.Error("删除文章类型失败",err)
		this.Redirect("/article/addType",302)
		return
	}
	o:=orm.NewOrm()

	var articleType models.ArticleType

	articleType.Id = id

	_,err=o.Delete(&articleType)
	if err!=nil{
		beego.Error("删除文章失败",err)
		this.Redirect("/article/addType",302)
		return
	}

	this.Redirect("/article/addType",302)

}
