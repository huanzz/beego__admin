package admin

import (
	. "bgadmin/common"
	m "bgadmin/models/admin"
	"errors"
	"github.com/astaxie/beego"
	"strings"
)


type NestPreparer interface {
	NestPrepare()
}

type BaseController struct {
	beego.Controller
	IsLogin  		bool			//标识 用户是否登陆
	Member    		m.Member 		//登陆的用户
}

// Check Login
func CheckLogin(membername string, password string) (member m.Member, err error) {
	member,_ = m.GetMemberByMemberName(membername)
	if member.Id == 0 {
		return member, errors.New("用户不存在")
	}
	if member.Password != PwdHash(password) {
		return member, errors.New("密码错误")
	}
	return member, nil
}

// Preparation
func (this *BaseController) Prepare() {
	this.IsLogin = false
	tu := this.GetSession("member_info")
	if tu != nil {
		if mber, ok := tu.(m.Member); ok {
			this.Member = mber
			this.Data["Member"] = mber
			this.IsLogin = true
		}
	}
	url := this.Ctx.Request.RequestURI
	menuMap := m.GetMenuMap(this.Member.Id)
	menuShortcut,_ := m.GetMenuShortcut(this.Member.Id)
	this.Data["menuShortcut"] = menuShortcut
	this.Data["WebTitle"] = menuMap[GetFullUrl(url)]
	this.Data["Member"] = this.Member
	this.Data["MenuView"] = m.GetMenuView(this.Member.Id)
	this.Layout = "admin/common/layout.html"

	// 判断子类是否实现了NestPreparer接口，如果实现了就调用接口方法。
	if app, ok := this.AppController.(NestPreparer); ok {
		app.NestPrepare()
	}

}

// no auth tips
func (this *BaseController)Tips()  {
	this.NoteAndJump("Warning", "没有操作权限", "/admin/index")
}

// Note And Jump
func (this *BaseController)NoteAndJump(types, msg, url string)  {
	this.Data["types"] = types
	this.Data["msg"] = msg
	this.Data["url"] = url
	this.TplName = "admin/common/jump.html"
	return
}

// Check Auth
func (this *BaseController)CheckAuth(memberId int,url string) bool  {
	if memberId !=1 {
		menuMap := m.GetMenuMap(this.Member.Id)
		webTitle := strings.Trim(menuMap[GetFullUrl(url)]," ")
		if webTitle == "" {
			return false
		}
	}
	return true
}


