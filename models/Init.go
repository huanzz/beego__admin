package models

import (
	"github.com/huanzz/bgadmin/common"
	"database/sql"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"log"
	_ "github.com/go-sql-driver/mysql"
	mo "github.com/huanzz/bgadmin/models/admin"
)

var o orm.Ormer


func SyncDataBase()  {
	CreateDb()
	Connect()
	o = orm.NewOrm()
	// 数据库别名
	name := "default"
	// drop table 后再建表
	force := true
	// 打印执行过程
	verbose := true
	// 遇到错误立即返回
	err := orm.RunSyncdb(name, force, verbose)
	if err != nil {
		fmt.Println(err)
	}

	 InsertInitMenu()		//Add Init Menu
	 InsertInitAuthGroup()	//Add Init	AuthGroup
	 InsertInitMember()		//Add Init Member

	fmt.Println("database init is complete.\nPlease restart the application")
}

func Connect()  {
	var dsn string
	db_type := beego.AppConfig.String("db_type")
	db_host := beego.AppConfig.String("db_host")
	db_port := beego.AppConfig.String("db_port")
	db_user := beego.AppConfig.String("db_user")
	db_pass := beego.AppConfig.String("db_pass")
	db_name := beego.AppConfig.String("db_name")
	if db_type == "mysql" {
		orm.RegisterDriver("mysql", orm.DRMySQL)
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&loc=Local", db_user, db_pass, db_host, db_port, db_name)
	}else {
		logs.Critical("Database driver is not allowed:", db_type)
	}
	orm.RegisterDataBase("default", db_type, dsn)
}

func CreateDb()  {
	db_type := beego.AppConfig.String("db_type")
	db_host := beego.AppConfig.String("db_host")
	db_port := beego.AppConfig.String("db_port")
	db_user := beego.AppConfig.String("db_user")
	db_pass := beego.AppConfig.String("db_pass")
	db_name := beego.AppConfig.String("db_name")

	var dsn string
	var sqlstring string
	if db_type == "mysql" {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8", db_user, db_pass, db_host, db_port)
		sqlstring = fmt.Sprintf("CREATE DATABASE  if not exists `%s` CHARSET utf8 COLLATE utf8_general_ci", db_name)
	}else {
		logs.Critical("Database driver is not allowed:", db_type)
	}
	db, err := sql.Open(db_type, dsn)
	if err != nil {
		panic(err.Error())
	}
	r, err := db.Exec(sqlstring)
	if err != nil {
		log.Println(err)
		log.Println(r)
	} else {
		log.Println("Database ", db_name, " created")
	}
	defer db.Close()
}

func InsertInitAuthGroup()  {
	InsertAuthGroup("Developer", "开发者", -2, 0, "1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25")
	InsertAuthGroup("SuperAdmin", "超级管理员",-2, 1,"1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25")
	InsertAuthGroup("Admin", "管理员",-2, 2, "1,2,3,4,5,6,7,8,9,10,11,12,13,15,16,17,18,20,21,22,23,24")
	InsertAuthGroup("User", "普通用户",-2, 2, "1,2,3,4,5,6")

}

func InsertAuthGroup(name string, des string, status int,memberId int, rules string)  {
	fmt.Println("insert auth group "+ name +" ...")
	a := new(mo.AuthGroup)
	a.Name = name
	a.Describe = des
	a.Status = status
	a.MemberId = memberId
	a.Rules = rules
	o = orm.NewOrm()
	o.Insert(a)
	fmt.Println("insert auth group "+ name +" end")
}

func InsertInitMember()  {
	InsertMember("developer","developer",common.PwdHash("102030"),"admin@admin.com","13111111111", -2,0, 1)
	InsertMember("super_administrator","super",common.PwdHash("123456"),"admin@admin.com","13111111111", -2,1, 2)
	InsertMember("administrator","admin",common.PwdHash("111111"),"admin@admin.com","13111111111", -2,2, 3)
	InsertMember("TestUser","testuser",common.PwdHash("111111"),"admin@admin.com","13111111112", -2,2, 4)
}

func InsertMember(nick, name, pwd, email, phone string, status,parentId int, groupId int)  {
	fmt.Println("insert member "+ name +" ...")
	m := new(mo.Member)
	m.NickName = nick
	m.MemberName = name
	m.Password = pwd
	m.Email = email
	m.Mobile = phone
	m.ParentId = 0
	m.Status = status
	m.ParentId = parentId
	group := mo.GetAuthGroupById(groupId)
	m.AuthGroup = &group
	o = orm.NewOrm()
	o.Insert(m)
	fmt.Println("insert member "+ name +" end")
}

func InsertInitMenu()  {
	notHide := 1
	isHide := 2
	InsertMenu(0,"系统首页", "admin", "/index", "fa fa-tachometer",notHide)
	InsertMenu(0,"基本操作", "admin", "/operate", "fa fa-tachometer",isHide)
	oper := mo.GetMenuByName("基本操作")
	InsertMenu(oper.Id,"权限警告", "admin", "/operate/jump", "fa fa-warning",isHide)
	InsertMenu(oper.Id,"跳转提示", "admin", "/operate/tips", "fa fa-info",isHide)
	InsertMenu(oper.Id,"个人信息", "admin", "/operate/person", "fa fa-user-secret",isHide)
	InsertMenu(oper.Id,"修改信息", "admin", "/operate/savemsg", "fa fa-send",isHide)

	InsertMenu(0,"会员管理", "admin", "/member/index", "fa fa-users",notHide)
	InsertMenu(0,"权限管理", "admin", "/auth/index", "fa fa-key",notHide)
	InsertMenu(0,"菜单管理", "admin", "/menu/index", "fa fa-th-large",notHide)
	member := mo.GetMenuByName("会员管理")
	InsertMenu(member.Id,"会员列表", "admin", "/member/list", "fa fa-user",notHide)
	InsertMenu(member.Id,"会员新建", "admin", "/member/add", "fa fa-user-plus",notHide)
	InsertMenu(member.Id,"会员编辑", "admin", "/member/edit", "fa fa-edit", isHide)
	InsertMenu(member.Id,"会员保存", "admin", "/member/doedit", "fa fa-edit",isHide)
	InsertMenu(member.Id,"会员删除", "admin", "/member/del", "fa fa-user-times",isHide)
	authGroup := mo.GetMenuByName("权限管理")
	InsertMenu(authGroup.Id,"权限组列表", "admin", "/auth/list", "fa fa-unlock-alt", notHide)
	InsertMenu(authGroup.Id,"权限组添加", "admin", "/auth/add", "fa fa-plus-square",notHide)
	InsertMenu(authGroup.Id,"权限组编辑", "admin", "/auth/edit", "fa fa-pencil-square",isHide)
	InsertMenu(authGroup.Id,"权限组保存", "admin", "/auth/doedit", "fa fa-pencil-square",isHide)
	InsertMenu(authGroup.Id,"权限组删除", "admin", "/auth/del", "fa fa-minus-square",isHide)
	InsertMenu(authGroup.Id,"权限授权", "admin", "/auth/authorize", "fa fa-child",isHide)
	menu := mo.GetMenuByName("菜单管理")
	InsertMenu(menu.Id,"菜单列表", "admin", "/menu/list", "fa fa-list-alt",notHide)
	InsertMenu(menu.Id,"菜单新建", "admin", "/menu/add", "fa fa-plus-circle",notHide)
	InsertMenu(menu.Id,"菜单编辑", "admin", "/menu/edit", "fa fa-pencil",isHide)
	InsertMenu(menu.Id,"菜单保存", "admin", "/menu/doedit", "fa fa-pencil",isHide)
	InsertMenu(menu.Id,"菜单删除", "admin", "/menu/del", "fa fa-minus-circle",isHide)
}

func InsertMenu(pid int, name, module, url, icon string, hide int)  {
	fmt.Println("insert Menu "+ url +" ...")
	m := new(mo.Menu)
	m.Name = name
	m.Pid = pid
	m.Module = module
	m.Url = url
	m.Icon = icon
	m.Sort = 1
	m.IsHide = hide
	m.IsShortcut = 1
	m.Status = -2
	o = orm.NewOrm()
	o.Insert(m)
	fmt.Println("insert  Menu "+ url +" end")
}

