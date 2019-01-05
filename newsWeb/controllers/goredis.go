package controllers

import (
	"github.com/astaxie/beego"
	"github.com/gomodule/redigo/redis"
)

type GoRedis struct {
	beego.Controller
}

func (this *GoRedis)ShowGet()  {
	//
	conn,err:=redis.Dial("tcp","127.0.0.1:6379")
	if err!=nil{
		beego.Error("连接redis失败",err)
	}
	reply,err:=conn.Do("set","11111","11111")
	replyString,err:=redis.String(reply,err)
	if err!=nil{
		beego.Error("设置数据失败",err)
	}
	beego.Info(replyString,"设置数据成功")
	reply,err=conn.Do("get","11111")
	replyString,err=redis.String(reply,err)
	if err!=nil{
		beego.Error("设置数据失败",err)
	}
	beego.Info(replyString,"设置数据成功")


	//
	resp,err:=conn.Do("mset","k1","v1","k2","v2")
	num,err:=redis.String(resp,err)
	if err!=nil{
		beego.Error("设置数据失败",err)
	}
	beego.Info("设置数据成功",num)
	reply,err=conn.Do("mget","k1","k2")
	replys,err:=redis.Values(reply,err)
	if err!=nil{
		beego.Error("获取数据失败",err)
	}
	var string1 string
	var string2 string
	redis.Scan(replys,&string1,&string2)
	beego.Info("replys",reply)
	this.Ctx.WriteString(string1+string2)
	//var buffer bytes.Buffer
	//enc:=gob.NewEncoder(&buffer)
	//enc.Encode(&article)
	//resp,err=conn.Do("get",buffer.Bytes())
	//
	//
	//resp,err=conn.Do("get","12")
	//res,err:=redis.Bytes(resp,err)
	//dec:=gob.NewDecoder(bytes.NewReader(res))
	//var testTyps []models.ArticleType
	//dec.Decode(&testTyps)
}
