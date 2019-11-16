package controllers

import (
	"bgadmin/common"
	"bgadmin/controllers/admin"
	"errors"
	"github.com/astaxie/beego"
	mo "bgadmin/models/admin"
)

type LoginController struct {
	beego.Controller
}


func (this *LoginController)Login() {
	memberInfo := this.GetSession("member_info")
	if memberInfo != nil {
		this.Ctx.Redirect(302, "/admin/index")
	}else {
		this.TplName = "login.html"
	}
}

func (this *LoginController)LoginAdmin() {
	memberName := this.GetString("memberName")
	password := this.GetString("password")
	member, err := admin.CheckLogin(memberName, password)
	if err == nil {
		this.SetSession("member_info", member)
		this.Notice("Success","登录成功","/admin/index")
	}else {
		this.Notice("Warning","登录失败，请检查账号或密码是否正确","/login")
	}

}

//退出
func (this *LoginController) Logout() {
	this.DelSession("member_info")
	this.Ctx.Redirect(302, "/login")
}

func (this *LoginController) ChangePwd()  {
	this.Data["WebTitle"] = "修改密码"
	this.TplName = "register.html"
}

func (this *LoginController)SaveChange()  {
	memberName := this.GetString("MemberName")
	Email := this.GetString("Email")
	Password := this.GetString("Password")
	RePassword := this.GetString("RePassword")
	var err error
	var member mo.Member
	member,err = mo.GetMemberByMemberName(memberName)
	if err == nil {
		if Email == member.Email {
			if Password == RePassword && Password != "" {
				member.Password = common.PwdHash(Password)
				_,err = mo.UpdateMember(&member)
			}else {
				err = errors.New("两次密码不一致！")
			}
		}else {
			err = errors.New("邮箱错误！")
		}
	}else {
		err = errors.New("用户名错误！")
	}

	if err == nil {
		this.Notice("Success","修改成功！","/login")
	}else {
		this.Notice("Error","修改失败:"+err.Error(),"/changepwd")
	}
}

func (this *LoginController)Notice(types, msg, url string)  {
	this.Data["types"] = types
	this.Data["msg"] = msg
	this.Data["url"] = url
	this.TplName = "notice.html"
	return
}
