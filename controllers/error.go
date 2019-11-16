package controllers

import "github.com/astaxie/beego"

type ErrorControl struct {
	beego.Controller
}

func (this *ErrorControl) Error404() {
	this.TplName = "error/404.html"
}

