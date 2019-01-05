package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type User struct {
	Id int
	UserName string `orm:"unique"`
	Pwd string
	Articles []*Article `orm:"rel(m2m)"`//多对多(可互换)
}

type Article struct {
	Id int `orm:"pk;auto"`
	Title string `orm:"size(100)"`
	Time time.Time `orm:"type(datetime);auto_now"`
	Count int `orm:"default(0)"`
	Img string `orm:"null"`
	Content string
	price float64 `orm:"digits(10);decimals(2)"`
	ArticleType *ArticleType `orm:"rel(fk)"` //1对多关系 多的
	Users []*User `orm:"reverse(many)"`//多对多(可互换)
}
//需要创建一对多的类型表以及多对多的用户与文章关系表
type ArticleType struct {
	Id int
	TypeName string `orm:"unique"`
	Article []*Article `orm:"reverse(many)"`
}


func init(){
	//建表的三步骤
	//注册数据库
	//第一个,为什么要用别名
	orm.RegisterDataBase("default","mysql","root:wujiu59@tcp(127.0.0.1:3306)/newsWeb?charset=utf8")
	//注册表
	orm.RegisterModel(new(User),new(Article),new(ArticleType))
	//跑起来
	orm.RunSyncdb("default",false,true)
}
