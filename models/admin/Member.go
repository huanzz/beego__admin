package admin

import (
	. "github.com/huanzz/bgadmin/common"
	"errors"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"log"
	"strings"
	"time"
)

type Member struct {
	Id			int
	NickName 	string		`orm:"unique" valid:"Required" form:"NickName"`
	MemberName 	string		`orm:"unique" valid:"Required" form:"MemberName"`
	Password 	string		`valid:"Required;MaxSize(20);MinSize(6)" `
	RePassword 	string		`orm:"-" valid:"Required;MaxSize(20);MinSize(6)""`
	Email 		string		`valid:"Email" form:"Email"`
	Mobile		string		`valid:"Mobile" form:"Mobile"`
	Status 		int			`orm:"default(1)"  form:"Status"`
	ParentId 	int			`orm:"default(0)"`
	IsShare 	int			`orm:"default(1)"`
	IsInside 	int			`orm:"default(1)"`
	UpdateAt	time.Time	`orm:"auto_now;type(datetime)"`
	CreateAt 	time.Time	`orm:"type(datetime);auto_now_add"`
	AuthGroup	*AuthGroup	`orm:"rel(fk)"`
}

type MemberL struct {
	Member
	Name string
}

// Register Model Member
func init()  {
	orm.RegisterModel(new(Member))
}

// check Password and RePassword
func (u *Member) Valid(v *validation.Validation) {
	if u.Password != u.RePassword {
		v.SetError("RePassword", "两次输入的密码不一样")
	}
}

// Check Member
func CheckMember(m *Member) (err error) {
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

// Insert Member
func InsertMember(m *Member) (int64, error) {
	if err := CheckMember(m); err!=nil {
		return 0, err
	}
	member := new(Member)
	member.NickName = m.NickName
	member.MemberName = m.MemberName
	member.Password = PwdHash("123456")
	member.Email = m.Email
	member.Mobile = m.Mobile
	member.ParentId = m.ParentId
	member.IsShare = m.IsShare
	member.IsInside = m.IsInside
	member.Status = m.Status
	member.AuthGroup = m.AuthGroup
	o := orm.NewOrm()
	id, err := o.Insert(member)
	return id, err
}

// Update Member message
func UpdateMember(m *Member) (int64, error) {
	if err := CheckMember(m); err != nil {
		return 0, err
	}
	member := GetMemberById(m.Id)
	if len(m.NickName) > 0 {
		member.NickName = m.NickName
	}
	if len(m.MemberName) > 0 {
		member.MemberName= m.MemberName
	}
	if len(m.Password) > 0 {
		member.Password = PwdHash(m.Password)
	}
	if len(m.Email) > 0 {
		member.Email = m.Email
	}
	if len(m.Mobile) > 0 {
		member.Mobile = m.Mobile
	}
	if m.ParentId != -1 {
		member.ParentId = m.ParentId
	}
	if m.IsShare != -1 {
		member.IsShare = m.IsShare
	}
	if m.IsInside != -1 {
		member.IsInside = m.IsInside
	}
	if m.Status != -1 {
		member.Status = m.Status
	}
	if m.AuthGroup != nil {
		member.AuthGroup = m.AuthGroup
	}
	o := orm.NewOrm()
	num, err := o.Update(&member)
	return num, err
}

// Delete Member By Id
func DelMemberById(id int) (int64, error) {
	o := orm.NewOrm()
	status, err := o.Delete(&Member{Id:id})
	return status,err
}

// Get Member By Member Name
func GetMemberByMemberName(membername string) (member Member, err error) {
	member = Member{MemberName: membername}
	o := orm.NewOrm()
	err = o.Read(&member, "MemberName")
	return member,err
}

// Get Member By Id
func GetMemberById(memberId int) (member Member) {
	member = Member{Id: memberId}
	o := orm.NewOrm()
	o.Read(&member)
	return member
}

// Exist Member By GroupId?
func ExistMemberByGroupId(groupId int) (res bool)  {
	res = false
	o := orm.NewOrm()
	db := o.QueryTable(new(Member))
	count,_ := db.Filter("AuthGroup__Id",groupId).Count()
	if count > 0{
		res = true
	}
	return res
}

// Get Member List using SQL
func GetMemberListInSQL(page int, pageSize int, sort string, search string, parentId int) (members []MemberL, count int64) {
	offset := PageOffset(page, pageSize)
	qb, _ := orm.NewQueryBuilder("mysql")
	var conn string
	if parentId ==1 {
		conn = "member.parent_id <> 0 AND (member.nick_name LIKE ? OR member.member_name LIKE ? OR member.mobile LIKE ? OR member.email LIKE ?)"
	}else {
		conn = "member.parent_id = ? AND (member.nick_name LIKE ? OR member.member_name LIKE ? OR member.mobile LIKE ? OR member.email LIKE ?)"
	}
	qb.Select("member.id","member.member_name","member.nick_name","member.email","member.mobile",
		"member.`status`","member.parent_id","member.is_share","member.is_inside","auth_group.`name`").
		From("member").
		InnerJoin("auth_group").On("member.auth_group_id = auth_group.id").
		Where(conn).
		OrderBy("member."+sort).Asc().Limit(pageSize).Offset(offset)
	sql := qb.String()
	o := orm.NewOrm()
	str := "%" + strings.Trim(search, " ") + "%"
	if parentId == 1 {
		count, _ = o.Raw(sql,str,str,str,str).QueryRows(&members)
	}else {
		count, _ = o.Raw(sql,parentId,str,str,str,str).QueryRows(&members)
	}
	return members,count
}