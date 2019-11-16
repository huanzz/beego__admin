package routers

import (
	"bgadmin/controllers"
	"bgadmin/controllers/admin"
	"bgadmin/filter"
	"github.com/astaxie/beego"
)

func init() {
	beego.ErrorController(&controllers.ErrorControl{})
    beego.Router("/", &controllers.MainController{})
	beego.Router("/login", &controllers.LoginController{},"*:Login")
	beego.Router("/logout", &controllers.LoginController{},"*:Logout")
	beego.Router("/loginadmin", &controllers.LoginController{},"*:LoginAdmin")
	beego.Router("/changepwd", &controllers.LoginController{},"*:ChangePwd")
	beego.Router("/savechange", &controllers.LoginController{},"*:SaveChange")
	beego.Router("/notice", &controllers.LoginController{},"*:Notice")

	beego.InsertFilter("/admin/*",beego.BeforeRouter,filter.FilterAdmin)
	beego.Router("/admin/operate/tips", &admin.BaseController{},"*:Tips")
	beego.Router("/admin/operate/jump", &admin.BaseController{},"*:NoteAndJump")
	beego.Router("/admin/operate/person", &admin.AdminController{},"*:Person")
	beego.Router("/admin/operate/savemsg", &admin.AdminController{},"*:SaveMsg")
	beego.Router("/admin/index", &admin.AdminController{},"*:Index")
	beego.Router("/admin/index/index", &admin.AdminController{},"*:Index")

	beego.Router("/admin/member/list", &admin.AdminController{},"*:MemberList")
	beego.Router("/admin/member/add", &admin.AdminController{},"*:MemberAddPage")
	beego.Router("/admin/member/edit", &admin.AdminController{},"*:MemberEditPage")
	beego.Router("/admin/member/doedit", &admin.AdminController{},"*:MemberEdit")
	beego.Router("/admin/member/del", &admin.AdminController{},"*:MemberDel")

	beego.Router("/admin/auth/list", &admin.AdminController{},"*:GroupList")
	beego.Router("/admin/auth/add", &admin.AdminController{},"*:GroupAddPage")
	beego.Router("/admin/auth/edit", &admin.AdminController{},"*:GroupEditPage")
	beego.Router("/admin/auth/doedit", &admin.AdminController{},"*:GroupEdit")
	beego.Router("/admin/auth/del", &admin.AdminController{},"*:GroupDel")
	beego.Router("/admin/auth/authorize", &admin.AdminController{},"*:AuthorizePage")
	beego.Router("/admin/auth/authorizeto", &admin.AdminController{},"*:Authorize")

	beego.Router("/admin/menu/list", &admin.AdminController{},"*:MenuList")
	beego.Router("/admin/menu/add", &admin.AdminController{},"*:MenuAddPage")
	beego.Router("/admin/menu/edit", &admin.AdminController{},"*:MenuEditPage")
	beego.Router("/admin/menu/doedit", &admin.AdminController{},"*:MenuEdit")
	beego.Router("/admin/menu/del", &admin.AdminController{},"*:MenuDel")
}
