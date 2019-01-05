package main

import "github.com/astaxie/beego"

import _"newsWeb/routers"

import _"newsWeb/models"

func main() {
	beego.AddFuncMap("prev",getPrev)
	beego.AddFuncMap("next",getNext)
	beego.Run()
}

func getPrev(pageIndex int)int  {
	 if pageIndex-1 <1{
	 	return 1
	 }
	 return pageIndex-1
}

func getNext(pageIndex int,pageCount int)int  {
	if pageIndex+1>pageCount{
		return pageCount
	}
	return pageIndex+1
}
/*
1.请求
2.路由
3.控制器
4.返回数据和视图
*/

/*
1.获取数据
2.校验数据
3.处理数据
4.返回数据
*/
