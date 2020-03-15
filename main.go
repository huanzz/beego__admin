package main

import (
	"bgadmin/models"
	"bgadmin/models/admin"
	_ "bgadmin/routers"
	"encoding/gob"
	"github.com/astaxie/beego"
	"os"
)

func main() {
	initialize()
	initSession()
	beego.Run()
}

func initialize()  {
	initArgs()
	models.Connect()
}


func initArgs() {
	args := os.Args
	for _, v := range args {
		if v == "syncdb" {
			models.SyncDataBase()
			os.Exit(0)
		}
	}
}


func initSession()  {
	//beego的session序列号是用gob的方式，因此需要将注册m.Member
	gob.Register(admin.Member{})
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.BConfig.WebConfig.Session.SessionName = "admin-key"
	beego.BConfig.WebConfig.Session.SessionProvider = "file"
	beego.BConfig.WebConfig.Session.SessionProviderConfig = "data/session"
}
