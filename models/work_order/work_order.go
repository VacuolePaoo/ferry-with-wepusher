type WorkOrderInfo struct {
	base.Model
	Title         string `gorm:"column:title; type:varchar(128)" json:"title" form:"title"`                    // 工单标题
	ProcessId     int    `gorm:"column:process_id; type:int(11)" json:"process_id" form:"process_id"`          // 流程ID
	Classify      int    `gorm:"column:classify; type:int(11)" json:"classify" form:"classify"`                // 分类ID
	State         string `gorm:"column:state; type:varchar(32)" json:"state" form:"state"`                     // 工单状态
	Priority      string `gorm:"column:priority; type:varchar(32)" json:"priority" form:"priority"`            // 工单优先级
	IsEnd         int    `gorm:"column:is_end; type:int(11)" json:"is_end" form:"is_end"`                      // 是否结束
	IsDenied      int    `gorm:"column:is_denied; type:int(11)" json:"is_denied" form:"is_denied"`             // 是否拒绝
	FormData      string `gorm:"column:form_data; type:longtext" json:"form_data" form:"form_data"`            // 表单数据
	Creator       int    `gorm:"column:creator; type:int(11)" json:"creator" form:"creator"`                    // 创建人
	RelatedPerson string `gorm:"column:related_person; type:varchar(255)" json:"related_person" form:"related_person"` // 相关人
	WechatOpenid  string `gorm:"column:wechat_openid; type:varchar(64)" json:"wechat_openid" form:"wechat_openid"`    // 微信openid
} 