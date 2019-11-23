package controllers

import "github.com/astaxie/beego"

type ErrorControl struct {
	beego.Controller
}

//	404 page
func (this *ErrorControl) Error404() {
	this.TplName = "error/404.html"
}

