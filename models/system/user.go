package system

import (
	"ferry/models/base"
)

/*
  @Author : lanyulei
*/

type SysUser struct {
	base.Model
	Username    string `gorm:"column:username; type:varchar(64)" json:"username" form:"username"`        // 用户名
	Password    string `gorm:"column:password; type:varchar(128)" json:"password" form:"password"`       // 密码
	NickName    string `gorm:"column:nick_name; type:varchar(64)" json:"nick_name" form:"nick_name"`     // 昵称
	Phone       string `gorm:"column:phone; type:varchar(11)" json:"phone" form:"phone"`                 // 手机号
	RoleId      int    `gorm:"column:role_id; type:int(11)" json:"role_id" form:"role_id"`              // 角色ID
	Avatar      string `gorm:"column:avatar; type:varchar(255)" json:"avatar" form:"avatar"`             // 头像
	Sex         string `gorm:"column:sex; type:varchar(255)" json:"sex" form:"sex"`                      // 性别
	Email       string `gorm:"column:email; type:varchar(128)" json:"email" form:"email"`                // 邮箱
	DeptId      int    `gorm:"column:dept_id; type:int(11)" json:"dept_id" form:"dept_id"`              // 部门ID
	PostId      int    `gorm:"column:post_id; type:int(11)" json:"post_id" form:"post_id"`              // 岗位ID
	Remark      string `gorm:"column:remark; type:varchar(255)" json:"remark" form:"remark"`             // 备注
	Status      string `gorm:"column:status; type:varchar(4)" json:"status" form:"status"`               // 状态
	WechatOpenid string `gorm:"column:wechat_openid; type:varchar(64)" json:"wechat_openid" form:"wechat_openid"` // 微信openid
}

func (SysUser) TableName() string {
	return "sys_user"
} 