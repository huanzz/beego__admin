package admin

import (
	"bgadmin/common"
	"bgadmin/models/admin"
	"errors"
	"strconv"
	"strings"
)

type AdminController struct {
	BaseController
}

//--------------- 首页及操作页 ----------------------------

//后台首页
func (this *AdminController)Index()  {
	this.Data["Member"] = this.Member
	this.TplName = "admin/index.html"
}

//个人信息页面
func (this *AdminController)Person()  {
	Id,_ := this.GetInt("Id")
	var user admin.Member
	user = admin.GetMemberById(Id)
	this.Data["User"] = user
	this.TplName = "admin/person.html"
}

//保存个人信息
func (this *AdminController)SaveMsg()  {
	Id,_ := this.GetInt("Id")
	oldPwd := this.GetString("OldPassword")
	newPwd := this.GetString("NewPassword")
	reNewPwd := this.GetString("ReNewPassword")
	mber := admin.Member{}
	var err error
	if err := this.ParseForm(&mber); err != nil {
		err = errors.New("请稍后再试")
	}
	if oldPwd == "" ||  newPwd == "" || newPwd == ""{
		_,err = admin.UpdateMember(&mber)
	}else {
		if newPwd != reNewPwd  {
			err = errors.New("新密码不一致")
		}else{
			member := admin.GetMemberById(Id)
			if member.Password != common.PwdHash(oldPwd) {
				err = errors.New("旧密码错误")
			}else {
				mber.Password = newPwd
				_,err = admin.UpdateMember(&mber)
			}
		}
	}
	if err == nil {
		this.NoteAndJump("Success","修改个人信息成功,","/admin/index")
	}else {
		this.NoteAndJump("Error","修改个人信息失败,"+err.Error(),"/admin/operate/person?Id="+strconv.Itoa(Id))
	}

}

//--------------- 会员管理 ----------------------------

//会员列表页
func (this *AdminController)MemberList()  {
	search := this.GetString("search")
	page,_ := this.GetInt("page")
	if page == 0{
		page = 1
	}
	members, total := admin.GetMemberListInSQL(page, 10, "id", search, this.Member.Id)
	this.Data["search"]= search
	this.Data["members"] = members
	this.Data["paginator"] = common.Paginator(page,10, total)
	this.TplName = "admin/member.html"
}

//会员新建页面
func (this *AdminController)MemberAddPage(){
	SelectGroup,count := admin.GetGroupList(1,100,"", this.Member.Id)
	if count > 0 {
		user := admin.Member{Status:1}
		this.Data["SelectGroup"] = SelectGroup
		this.Data["User"] = user
		this.TplName = "admin/memberEdit.html"
	} else{
		this.NoteAndJump("Warning","操作失败，请添加权限组后再试","/admin/member/list")
	}
}

//会员编辑页
func (this *AdminController)MemberEditPage(){
	Id,_ := this.GetInt("Id")
	var user admin.Member
	user = admin.GetMemberById(Id)
	SelectGroup,_ := admin.GetGroupList(1,100,"",this.Member.Id)
	this.Data["SelectGroup"] = SelectGroup
	this.Data["User"] = user
	this.TplName = "admin/memberEdit.html"
}

//会员信息编辑操作
func (this *AdminController)MemberEdit(){
	authId,_ := this.GetInt("authGroup")
	member := admin.Member{}
	if err := this.ParseForm(&member); err != nil {
		this.NoteAndJump("Warning","获取会员信息失败,"+err.Error(),"/admin/member/list")
	}
	authGroup := admin.GetAuthGroupById(authId)
	member.AuthGroup= &authGroup
	var err error
	str := ""
	member.ParentId = this.Member.Id
	if member.Id > 0{
		_,err = admin.UpdateMember(&member)
		str = "更新"
	} else {
		_,err = admin.InsertMember(&member)
		str = "添加"
	}

	if err == nil {
		this.NoteAndJump("Success",str+"会员信息成功", "/admin/member/list")
	}else {
		this.NoteAndJump("Error",str+"会员信息失败", "/admin/member/list")
	}

}

//会员信息删除操作
func (this *AdminController)MemberDel(){
	Id,_ :=this.GetInt("Id")
	mber := admin.GetMemberById(Id)
	var err error
	if mber.Status == -2 {
		this.NoteAndJump("Warning","默认会员不可删除", "/admin/auth/list")
	}else {
		_,err = admin.DelMemberById(Id)
		if err == nil {
			this.NoteAndJump("Success","删除会员信息成功", "/admin/member/list")
		}else {
			this.NoteAndJump("Error","删除会员信息失败", "/admin/member/list")
		}
	}



}


//--------------- 权限管理 ----------------------------

//权限组列表
func (this *AdminController)GroupList()  {
	search := this.GetString("search")
	page,_ := this.GetInt("page")
	if page == 0{
		page = 1
	}
	authGroup,_ := admin.GetGroupInSQL(page, 10, search, this.Member.Id)
	total := admin.TotalAuthGroup(search, this.Member.Id)
	this.Data["authGroup"] = authGroup
	this.Data["search"] = search
	this.Data["paginator"] = common.Paginator(page,10, total)
	this.TplName = "admin/authgroup.html"
}

//权限组新建页面
func (this *AdminController)GroupAddPage(){
	initRules := "1,2,3,4,5,6"
	authGroup := admin.AuthGroup{Status:1,Rules:initRules}
	authParent := admin.GetAuthGroupById(this.Member.AuthGroup.Id)
	SelectGroup := admin.GetMenuSelect(initRules, authParent.Rules)
	this.Data["SelectGroup"] = SelectGroup
	this.Data["authGroup"] = authGroup
	this.Data["IsAdd"] = true
	this.TplName = "admin/authEdit.html"
}

//权限组编辑页面
func (this *AdminController)GroupEditPage(){
	Id,_ := this.GetInt("Id")
	this.Data["IsAdd"] = false
	var authGroup admin.AuthGroup
	authGroup = admin.GetAuthGroupById(Id)
	this.Data["authGroup"] = authGroup
	this.TplName = "admin/authEdit.html"
}

//权限组信息编辑操作
func (this *AdminController)GroupEdit(){
	authGroup := admin.AuthGroup{}
	if err := this.ParseForm(&authGroup); err != nil {
		this.NoteAndJump("Warning","获取权限组信息失败,"+err.Error(),"/admin/auth/list")
	}
	var err error
	str := ""
	if authGroup.Id > 0{
		_,err = admin.UpdateAuthGroup(&authGroup)
		str = "更新"
	} else {
		authGroup.Rules = "1,2,3,4,5,6"
		authGroup.MemberId = this.Member.Id
		_,err = admin.InsertAuthGroup(&authGroup)
		str = "添加"
	}
	if err == nil {
		this.NoteAndJump("Success",str+"权限组信息成功", "/admin/auth/list")
	}else {
		this.NoteAndJump("Error",str+"权限组信息失败", "/admin/auth/list")
	}

}

//权限组删除操作
func (this *AdminController)GroupDel(){
	Id,_ :=this.GetInt("Id")
	auth := admin.GetAuthGroupById(Id)
	if auth.Status == -2 {
		this.NoteAndJump("Warning","默认权限组不可删除", "/admin/auth/list")
	}else{
		hasMemer := admin.ExistMemberByGroupId(auth.Id)
		if hasMemer {
			this.NoteAndJump("Error","操作失败，请删除该权限组的所有会员后重试", "/admin/auth/list")
		}else {
			_,err := admin.DelAuthGroupById(Id)
			if err == nil {
				this.NoteAndJump("Success","删除权限组信息成功", "/admin/auth/list")
			}else {
				this.NoteAndJump("Error","删除权限组信息失败", "/admin/auth/list")
			}
		}
	}
}

//权限组授权页面
func (this *AdminController)AuthorizePage()  {
	Id,_ := this.GetInt("Id")
	auth := admin.GetAuthGroupById(Id)
	authParent := admin.GetAuthGroupById(this.Member.AuthGroup.Id)
	SelectGroup := admin.GetMenuSelect(auth.Rules, authParent.Rules)
	this.Data["SelectGroup"] = SelectGroup
	this.Data["Auth"] = auth
	this.TplName = "admin/authorize.html"
}

//权限组授权操作
func (this *AdminController) Authorize() {
	Id,_ := this.GetInt("Id")
	rules := this.GetStrings("rules")
	ids := strings.Join(rules,",")
	auth := admin.GetAuthGroupById(Id)
	auth.Rules = ids
	_,err := admin.UpdateAuthGroup(&auth)
	if err == nil {
		this.NoteAndJump("Success","权限修改成功","/admin/auth/list")
	}else {
		this.NoteAndJump("Error","权限修改失败,"+err.Error(),"/admin/auth/authorize?Id="+string(Id))
	}
}


//--------------- 菜单管理 ----------------------------

//菜单列表页面
func (this *AdminController)MenuList()  {
	Pid,_ := this.GetInt("Pid")
	search := this.GetString("search")
	page,_ := this.GetInt("page")
	if page == 0{
		page = 1
	}
	var pageSize int
	if Pid == 0 {
		pageSize = 10
	}else {
		pageSize = 50
	}
	memberId := this.Member.Id
	menuList,_ := admin.GetMenuListInSQL(page, pageSize, search, memberId, Pid)
	total := admin.TotalMenuList(search, memberId, Pid)
	this.Data["Pid"] = Pid
	this.Data["search"] = search
	this.Data["menuList"] = menuList
	this.Data["ParentMenu"] = admin.GetMenuById(Pid)
	this.Data["paginator"] = common.Paginator(page,pageSize, total)
	this.TplName = "admin/menu.html"

}

//菜单新建页面
func (this *AdminController)MenuAddPage(){
	menu := admin.Menu{Status:1,Module:"admin",Sort:1,IsHide:1,IsShortcut:1 }
	memberId := this.Member.Id
	selectMenu := admin.GetMenuView(memberId)
	this.Data["selectMenu"] = selectMenu
	this.Data["menu"] = menu
	this.TplName = "admin/menuEdit.html"
}

//菜单编辑页面
func (this *AdminController)MenuEditPage()  {
	Id,_ := this.GetInt("Id")
	var menu admin.Menu
	menu = admin.GetMenuById(Id)
	memberId := this.Member.Id
	selectMenu := admin.GetMenuView(memberId)
	this.Data["selectMenu"] = selectMenu
	this.Data["menu"] = menu
	this.TplName = "admin/menuEdit.html"
}

//菜单编辑操作
func (this *AdminController)MenuEdit(){
	menu := admin.Menu{}
	if err := this.ParseForm(&menu); err != nil {
		this.NoteAndJump("Warning","获取权限组信息失败,"+err.Error(),"/admin/menu/list")
	}
	var err error
	str := ""
	if menu.Id > 0{
		_,err = admin.UpdateMenu(&menu)
		str = "更新"
	} else {
		var id int64
		id,err = admin.InsertMenu(&menu)
		admin.AddRuleWhenAddMenu(id, this.Member.Id)
		str = "添加"
	}
	if err == nil {
		this.NoteAndJump("Success",str+"菜单信息成功", "/admin/menu/list")
	}else {
		this.NoteAndJump("Error",str+"菜单信息失败", "/admin/menu/list")
	}
}

//菜单删除操作
func (this *AdminController)MenuDel(){
	Id,_ :=this.GetInt("Id")
	menu := admin.GetMenuById(Id)
	if menu.Status == -2 {
		this.NoteAndJump("Warning","默认菜单不可删除", "/admin/menu/list")
	}else {
		_,err := admin.DelMenuById(Id)
		admin.DelOneRuleByMenu(Id, this.Member.Id)
		if err == nil {
			this.NoteAndJump("Success","删除菜单组信息成功", "/admin/menu/list")
		}else {
			this.NoteAndJump("Error","删除权菜单信息失败", "/admin/menu/list")
		}
	}


}

//
func (this *AdminController)MenuIcon()  {
	i,_  := this.GetInt("Id")
	this.Layout = ""
	if i == 1111 {
		this.TplName = "admin/common/icons.html"
	}else {
		this.TplName = "admin/common/empty.html"
	}
}
