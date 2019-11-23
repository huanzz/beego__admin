package admin

import (
	. "bgadmin/common"
	"errors"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"log"
	"strings"
	"time"
)

type Menu struct {
	Id			int
	Name 		string		`orm:"unique" valid:"Required" form:"Name"`
	Pid 		int			`orm:"default(0)" form:"Pid"`
	Sort 		int			`orm:"default(1)" form:"Sort"`
	Module	 	string		`valid:"Required" form:"Module"`
	Url		 	string		`valid:"Required" form:"Url"`
	IsHide 		int			`orm:"default(1)" form:"IsHide"`
	IsShortcut 	int			`orm:"default(1)" form:"IsShortcut"`
	Status 		int			`orm:"default(1)" form:"Status"`
	Icon		string		`orm:"null" form:"Icon"`
	UpdateAt	time.Time	`orm:"auto_now;type(datetime)"`
	CreateAt 	time.Time	`orm:"type(datetime);auto_now_add"`
	Active 		bool		`orm:"null"`
}

type MenuView struct {
	Menu
	Child 		[]Menu
	EmptyChild 	bool
	Active 		bool
}

type MenuList struct {
	Menu
	ChildCount 	int64
}

// Register Model Menu
func init()  {
	orm.RegisterModel(new(Menu))
}

// Check Menu Using valid
func CheckMenu(m *Menu) (err error) {
	valid := validation.Validation{}
	b,_ := valid.Valid(&m)
	if !b {
		for _, err := range valid.Errors {
			log.Println(err.Key, err.Message)
			return errors.New(err.Message)
		}
	}
	return nil
}

// Insert Menu
func InsertMenu(m *Menu) (int64, error) {
	if err := CheckMenu(m); err!=nil {
		return 0, err
	}
	menu := new(Menu)
	menu.Name = m.Name
	menu.Pid = m.Pid
	menu.Sort = m.Sort
	menu.Module = m.Module
	menu.Url = m.Url
	menu.IsHide = m.IsHide
	menu.IsShortcut = m.IsShortcut
	menu.Icon = m.Icon
	menu.Status = m.Status
	o := orm.NewOrm()
	id, err := o.Insert(menu)
	return id, err
}

// Update Menu
func UpdateMenu(m *Menu) (int64, error) {
	if err := CheckMenu(m); err != nil {
		return 0, err
	}
	menu := make(orm.Params)
	if len(m.Name) > 0 {
		menu["Name"] = m.Name
	}
	if len(m.Module) > 0 {
		menu["Module"] = m.Module
	}
	if len(m.Url) > 0 {
		menu["Url"] = m.Url
	}
	if len(m.Icon) > 0 {
		menu["Icon"] = m.Icon
	}
	if m.Pid != -1 {
		menu["Pid"] = m.Pid
	}
	if m.Sort != -1 {
		menu["Sort"] = m.Sort
	}
	if m.IsHide != -1 {
		menu["IsHide"] = m.IsHide
	}
	if m.IsShortcut != -1 {
		menu["IsShortcut"] = m.IsShortcut
	}
	if m.Status != -1 {
		menu["Status"] = m.Status
	}
	if len(menu) == 0{
		return 0, errors.New("update field is empty")
	}
	o := orm.NewOrm()
	var table Menu
	num, err := o.QueryTable(table).Filter("Id", m.Id).Update(menu)
	return num, err
}

// Delete Menu By Id
func DelMenuById(id int) (int64, error) {
	o := orm.NewOrm()
	status, err := o.Delete(&Menu{Id:id})
	return status,err
}

func GetMenuByName(name string) (menu Menu) {
	menu = Menu{Name: name}
	o := orm.NewOrm()
	o.Read(&menu,"Name")
	return menu
}

// Get one Menu By Id
func GetMenuById(Id int) (menu Menu) {
	menu = Menu{Id: Id}
	o := orm.NewOrm()
	o.Read(&menu)
	return menu
}

// Get Menu Ids by memberId
func GetMenuIds(memberId int) (ids []int) {
	member := GetMemberById(memberId)
	authGroup := GetAuthGroupById(member.AuthGroup.Id)
	rules := authGroup.Rules
	ids = StrToIntArr(rules)
	return ids
}

// Get MenuList generate View
func GetMenuView(memberId int) []MenuView {
	ids := GetMenuIds(memberId)
	var menu []Menu
	o := orm.NewOrm()
	db := o.QueryTable(new(Menu))
	db.Filter("Id__in", ids).OrderBy("sort").All(&menu)
	var res []MenuView
	for _,v := range menu {
		row := MenuView{Menu: v}
		db.Filter("Id__in", ids).Filter("pid", v.Id).OrderBy("sort").All(&row.Child)
		if len(row.Child) == 0{
			row.EmptyChild = true
		}
		res = append(res, row)
	}
	return res
}

// Get Menu Map (Menu-key-value)
func GetMenuMap(memberId int) (menuMap map[string]string) {
	ids := GetMenuIds(memberId)
	var menu []Menu
	o := orm.NewOrm()
	db := o.QueryTable(new(Menu))
	db.Filter("Id__in", ids).OrderBy("sort").All(&menu)
	menuMap = make(map[string]string)
	for _,v := range menu {
		str := "/"+v.Module + v.Url
		menuMap[str] = v.Name
	}
	return menuMap
}

// Get Menu Count By ParentId
func GetMenuCountByPid(pid int) (count int64)  {
	o := orm.NewOrm()
	db := o.QueryTable(new(Menu))
	count,_ = db.Filter("Pid", pid).Count()
	return count
}

// Get MenuList to generate Selection
func GetMenuSelect(rules string, parentRules string)(res []MenuView)  {
	var menu []Menu
	ids := StrToIntArr(rules)
	idsParent := StrToIntArr(parentRules)
	o := orm.NewOrm()
	db := o.QueryTable(new(Menu))
	db.Filter("Id__in", idsParent).OrderBy("sort").All(&menu)
	for _,v := range menu {
		row := MenuView{Menu:v}
		db.Filter("pid", v.Id).OrderBy("sort").All(&row.Child)
		if ok := NumInIds(v.Id, ids); ok{
			row.Active = true
		}
		if len(row.Child) == 0{
			row.EmptyChild = true
		}
		for _,val := range row.Child{
			if ok := NumInIds(val.Id, ids); ok{
				val.Active = true
			}
		}
		res = append(res, row)
	}
	return res
}

// Get Menu List of Shortcut
func GetMenuShortcut(memberId int) (res []Menu, count int64) {
	ids := GetMenuIds(memberId)
	o := orm.NewOrm()
	db := o.QueryTable(new(Menu))
	db.Filter("Id__in", ids).All(&res)
	count,_ = db.Filter("Id__in", ids).Count()
	return res, count
}

// Get Menu List Using SQL
func GetMenuListInSQL(page int, pageSize int, search string, memberId int, Pid int) (menuList []MenuList, count int64) {
	member := GetMemberById(memberId)
	authGroup := GetAuthGroupById(member.AuthGroup.Id)
	ids := authGroup.Rules

	offset := PageOffset(page, pageSize)
	qb, _ := orm.NewQueryBuilder("mysql")
	var menus []Menu
	o := orm.NewOrm()
	str := "%" + strings.Trim(search, " ") + "%"
	conn := "(menu.id in ("+ids+")) AND menu.pid = ? AND (menu.name LIKE ? OR menu.module LIKE ? OR menu.url LIKE ?)"
	qb.Select("menu.id","menu.name","menu.pid","menu.sort","menu.module","menu.url","menu.is_hide","menu.is_shortcut","menu.`status`").
		From("menu").
		Where(conn).
		OrderBy("menu.create_at").Asc().Limit(pageSize).Offset(offset)
	sql := qb.String()
	count, _ = o.Raw(sql,Pid,str,str,str).QueryRows(&menus)
	for _,v := range menus {
		row := MenuList{Menu: v}
		num := GetMenuCountByPid(v.Id)
		row.ChildCount = num
		menuList = append(menuList, row)
	}
	return menuList, count
}
