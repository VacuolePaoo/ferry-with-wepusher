func (b *BodyData) SendNotify() (err error) {
	var (
		emailList []string
		phoneList []string
	)

	switch b.Priority {
	case 1:
		b.PriorityValue = "正常"
	case 2:
		b.PriorityValue = "紧急"
	case 3:
		b.PriorityValue = "非常紧急"
	}

	for _, c := range b.Classify {
		switch c {
		case 1: // 邮件
			users := b.SendTo.(map[string]interface{})["userList"].([]system.SysUser)
			if len(users) > 0 {
				for _, user := range users {
					emailList = append(emailList, user.Email)
					phoneList = append(phoneList, user.Phone)
				}
				err = b.ParsingTemplate()
				if err != nil {
					logger.Errorf("模版内容解析失败，%v", err.Error())
					return
				}
				go email.SendMail(emailList, b.EmailCcTo, b.Subject, b.Content)
				dingtalkEnable := viper.GetBool("settings.dingtalk.enable")
				if dingtalkEnable {
					url := fmt.Sprintf("%s/#/process/handle-ticket?workOrderId=%d&processId=%d", b.Domain, b.Id, b.ProcessId)
					go dingtalk.SendDingMsg(phoneList, url, b.Title, b.Creator, b.PriorityValue, b.CreatedAt)
				}
			}
		case 2: // 微信
			wechatEnable := viper.GetBool("settings.wechat.enable")
			if wechatEnable {
				// 获取工单创建人openid
				var workOrder process.WorkOrderInfo
				err = orm.Eloquent.Model(&process.WorkOrderInfo{}).Where("id = ?", b.Id).Find(&workOrder).Error
				if err != nil {
					logger.Errorf("查询工单信息失败，%v", err.Error())
					return
				}
				if workOrder.CreatorOpenId != "" {
					url := fmt.Sprintf("%s/#/process/handle-ticket?workOrderId=%d&processId=%d", b.Domain, b.Id, b.ProcessId)
					go wechat.SendWorkOrderResult(
						workOrder.CreatorOpenId,
						b.Title,
						"通过",
						b.Description,
						b.CreatedAt,
					)
				}
			}
		}
	}
	return
}