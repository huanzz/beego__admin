package admin

import (
	. "bgadmin/common"
	"errors"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"log"
	"strconv"
	"strings"
	"time"
)

type AuthGroup struct {
	Id			int
	Name 		string		`orm:"unique" valid:"Required" form:"Name"`
	Describe 	string		`orm:"null" form:"Describe"`
	Status 		int			`orm:"default(1)" form:"Status"`
	Rules		string		`orm:"null"`
	MemberId	int			`orm:"null"`
	UpdateAt	time.Time	`orm:"auto_now;type(datetime)"`
	CreateAt 	time.Time	`orm:"type(datetime);auto_now_add"`
	Member 		[]*Member	`orm:"reverse(many)"`
}

// Register Model AuthGroup
func init()  {
	orm.RegisterModel(new(AuthGroup))
}

// Check AuthGroup Using valid
func CheckAuthGroup(a *AuthGroup) (err error) {
	valid := validation.Validation{}
	b,_ := valid.Valid(&a)
	if !b {
		for _, err := range valid.Errors {
			log.Println(err.Key, err.Message)
			return errors.New(err.Message)
		}
	}
	return nil
}

// Insert AuthGroup
func InsertAuthGroup(a *AuthGroup) (int64, error) {
	if err := CheckAuthGroup(a); err!=nil {
		return 0, err
	}
	auth := new(AuthGroup)
	auth.Name = a.Name
	auth.Describe = a.Describe
	auth.Rules = a.Rules
	auth.MemberId = a.MemberId
	auth.Status = a.Status
	o := orm.NewOrm()
	id, err := o.Insert(auth)
	return id, err
}

// Update AuthGroup
func UpdateAuthGroup(a *AuthGroup) (int64, error) {
	if err := CheckAuthGroup(a); err != nil {
		return 0, err
	}
	auth := make(orm.Params)
	if len(a.Name) > 0 {
		auth["Name"] = a.Name
	}
	if len(a.Describe) > 0 {
		auth["Describe"] = a.Describe
	}
	if len(a.Rules) > 0 {
		auth["Rules"] = a.Rules
	}
	if a.Status != -1 {
		auth["Status"] = a.Status
	}
	if a.MemberId != -1 {
		auth["MemberId"] = a.MemberId
	}
	if len(auth) == 0{
		return 0, errors.New("update field is empty")
	}
	o := orm.NewOrm()
	var table AuthGroup
	num, err := o.QueryTable(table).Filter("Id", a.Id).Update(auth)
	return num, err
}

// Delete AuthGroup By Id
func DelAuthGroupById(id int) (int64, error) {
	o := orm.NewOrm()
	status, err := o.Delete(&AuthGroup{Id: id})
	return status, err
}

// Get AuthGroup By Id
func GetAuthGroupById(id int) (ag AuthGroup) {
	ag = AuthGroup{Id: id}
	o := orm.NewOrm()
	o.Read(&ag, "Id")
	return ag
}

// Get Group List
func GetGroupList(page int, pageSize int, search string, memberId int) (authGroup []AuthGroup, count int64) {
	o := orm.NewOrm()
	db := o.QueryTable(new(AuthGroup))
	offset := PageOffset(page, pageSize)
	if memberId == 1{
		count,_ =db.Limit(pageSize, offset).Exclude("MemberId", 0).OrderBy("-CreateAt").All(&authGroup)
		//count, _ = db.Limit(pageSize, offset).Exclude("MemberId", 0).Count()
	}else {
		count, _ = db.Limit(pageSize, offset).Filter("MemberId", memberId).OrderBy("-CreateAt").All(&authGroup)
		//count, _ = db.Limit(pageSize, offset).Filter("MemberId", memberId).Count()
	}

	return authGroup, count
}

// Get AuthGroup Using SQL
func GetGroupInSQL(page int, pageSize int, search string, memberId int) (authGroup []AuthGroup, count int64) {
	offset := PageOffset(page, pageSize)
	qb, _ := orm.NewQueryBuilder("mysql")
	var conn string
	if memberId ==1 {
		conn = "auth_group.member_id <> 0 AND (auth_group.name LIKE ? OR auth_group.describe LIKE ? )"
	}else {
		conn = "auth_group.member_id = ? AND (auth_group.name LIKE ? OR auth_group.describe LIKE ?)"
	}
	qb.Select("auth_group.id","auth_group.name","auth_group.describe","auth_group.member_id","auth_group.`status`").
		From("auth_group").
		Where(conn).
		OrderBy("auth_group.create_at").Asc().Limit(pageSize).Offset(offset)
	sql := qb.String()
	o := orm.NewOrm()
	str := "%" + strings.Trim(search, " ") + "%"
	if memberId == 1 {
		count, _ = o.Raw(sql,str,str).QueryRows(&authGroup)
	}else {
		count, _ = o.Raw(sql,memberId,str,str).QueryRows(&authGroup)
	}
	return authGroup,count
}

// Total Auth Group
func TotalAuthGroup(search string, memberId int) (count int64)  {
	qb, _ := orm.NewQueryBuilder("mysql")
	var conn string
	if memberId ==1 {
		conn = "auth_group.member_id <> 0 AND (auth_group.name LIKE ? OR auth_group.describe LIKE ? )"
	}else {
		conn = "auth_group.member_id = ? AND (auth_group.name LIKE ? OR auth_group.describe LIKE ?)"
	}
	qb.Select("auth_group.id","auth_group.name","auth_group.describe","auth_group.member_id","auth_group.`status`").
		From("auth_group").
		Where(conn)
	sql := qb.String()
	o := orm.NewOrm()
	var authGroup []AuthGroup
	str := "%" + strings.Trim(search, " ") + "%"
	if memberId == 1 {
		count, _ = o.Raw(sql,str,str).QueryRows(&authGroup)
	}else {
		count, _ = o.Raw(sql,memberId,str,str).QueryRows(&authGroup)
	}
	return count
}

// Add Rules When Add Menu
func AddRuleWhenAddMenu(menuId int64, memberId int)  {
	member := GetMemberById(memberId)
	parentId := member.ParentId
	groupId := member.AuthGroup.Id
	auth := GetAuthGroupById(groupId)
	rules := auth.Rules
	rules = rules + "," + strconv.FormatInt(menuId, 10)
	auth.Rules = rules
	UpdateAuthGroup(&auth)

	for groupId != 1  {
		member = GetMemberById(parentId)
		groupId = member.AuthGroup.Id

		auth := GetAuthGroupById(groupId)
		rules := auth.Rules
		rules = rules + "," + strconv.FormatInt(menuId, 10)
		auth.Rules = rules
		UpdateAuthGroup(&auth)
	}

}


// Delete On eRule By Menu Id
func DelOneRuleByMenu(menuId int, memberId int)  {
	member := GetMemberById(memberId)
	parentId := member.ParentId
	groupId := member.AuthGroup.Id
	auth := GetAuthGroupById(groupId)
	rules := auth.Rules
	ids := GetMenuIds(memberId)
	var newIds []string
	for _,v := range ids {
		if v != menuId {
			newIds = append(newIds,strconv.Itoa(v))
		}
	}
	rules = strings.Join(newIds,",")
	auth.Rules = rules
	UpdateAuthGroup(&auth)

	for groupId != 1  {
		member = GetMemberById(parentId)
		groupId = member.AuthGroup.Id

		auth := GetAuthGroupById(groupId)
		rules := auth.Rules
		ids := GetMenuIds(memberId)
		var newIds []string
		for _,v := range ids {
			if v != menuId {
				newIds = append(newIds,strconv.Itoa(v))
			}
		}
		rules = strings.Join(newIds,",")
		auth.Rules = rules
		UpdateAuthGroup(&auth)

	}

}