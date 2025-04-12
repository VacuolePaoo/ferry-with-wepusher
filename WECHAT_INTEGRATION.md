# Ferry 工单系统 - 微信集成说明

本文档介绍了 Ferry 工单系统中集成的微信通知功能。该功能允许系统在工单处理完成后，通过微信模板消息向工单创建者发送通知。它使用微信网页授权获取用户的 OpenID 作为唯一标识，无需用户在工单系统内部注册账户。

## 功能概述

1.  **用户身份识别**: 通过微信网页授权 (`snsapi_base`) 获取用户的 OpenID。
2.  **OpenID 存储**: 在用户提交工单时，将其微信 OpenID 与工单关联存储。
3.  **工单完成通知**: 当工单流程结束时，向关联了 OpenID 的工单创建者发送微信模板消息通知处理结果。

## 配置说明

所有微信相关的配置都在 `config/settings.yml` 文件中的 `settings.wechat` 部分进行。

```yaml
settings:
  # ... 其他配置 ...
  wechat:
    # 是否启用微信通知功能 (true/false)
    enable: true

    # 微信公众号的 AppID
    appid: "你的微信公众号AppID"

    # 微信公众号的 AppSecret
    appsecret: "你的微信公众号AppSecret"

    # 用于发送工单结果通知的模板消息 ID
    # 需要在微信公众号后台申请和配置
    template_id: "你的模板消息ID"

    # 微信网页授权回调地址
    # 必须与微信公众号后台配置的网页授权域名下的路径一致
    # 例如: "http://your-ferry-domain.com/api/v1/wechat/callback"
    redirect_uri: "你的域名/api/v1/wechat/callback"

    # 模板消息格式定义
    # 定义了发送通知时消息各部分的内容和颜色
    # {变量名} 会被替换为实际值
    template_format:
      first: # 消息头部内容
        value: "您的工单已处理完成"
        color: "#173177" # 颜色代码，例如 #RRGGBB
      keyword1: # 对应模板中的 {{keyword1.DATA}}
        value: "{title}" # 使用 {title} 作为工单标题占位符
        color: "#173177"
      keyword2: # 对应模板中的 {{keyword2.DATA}}
        value: "{result}" # 使用 {result} 作为处理结果占位符 (通过/拒绝)
        color: "#173177"
      keyword3: # 对应模板中的 {{keyword3.DATA}}
        value: "{remarks}" # 使用 {remarks} 作为处理备注占位符
        color: "#173177"
      keyword4: # 对应模板中的 {{keyword4.DATA}}
        value: "{time}" # 使用 {time} 作为处理时间占位符
        color: "#173177"
      remark: # 消息尾部内容
        value: "感谢您的使用"
        color: "#173177"
```

**请务必将上述配置中的占位符替换为您自己的实际信息。**

## 微信公众号后台设置

为了使该功能正常工作，您需要在微信公众号管理后台进行以下配置：

1.  **网页授权域名**:
    *   在 "公众号设置" -> "功能设置" -> "网页授权域名" 中，添加 Ferry 系统部署的域名（与 `redirect_uri` 中的域名部分一致）。
2.  **模板消息**:
    *   在 "广告与服务" -> "模板消息" 中，申请并添加一个新的模板。
    *   模板的格式需要与 `config/settings.yml` 文件中的 `template_format` 相对应。例如，如果配置中有 `keyword1` 到 `keyword4`，则模板中也需要包含 `{{keyword1.DATA}}` 到 `{{keyword4.DATA}}`。
    *   获取该模板的 **模板 ID**，并填入 `config/settings.yml` 的 `template_id` 字段。
    *   确保模板消息功能已启用。
3.  **IP 白名单**:
    *   在 "基本配置" -> "IP白名单" 中，添加 Ferry 后端服务器的公网 IP 地址。否则，后端将无法调用微信 API 获取 Access Token 等信息。

## 用户使用流程

1.  **访问页面**: 用户通过微信或其他浏览器访问工单创建页面 (`/#/process/create-ticket`)。
2.  **检查授权**: 页面加载时，系统会检查浏览器 Cookie 中是否已存在用户的 OpenID。
3.  **微信授权**:
    *   如果 Cookie 中**没有** OpenID，页面右上角会显示 "微信授权" 按钮。
    *   用户点击按钮，页面会跳转到微信授权页面。
    *   用户同意授权。
4.  **获取 OpenID**: 微信将用户重定向回 `redirect_uri` 配置的地址 (`/api/v1/wechat/callback`)，并附带 `code` 参数。
5.  **后端处理**: 后端服务使用 `code`、`appid` 和 `appsecret` 向微信服务器换取用户的 `openid`。
6.  **存储 OpenID**: 后端将获取到的 `openid` 写入用户浏览器的 Cookie 中，并将用户重定向回工单创建页面 (`/#/process/create-ticket`)。
7.  **提交工单**: 用户填写工单信息并提交。前端代码会自动从 Cookie 读取 `openid` 并将其包含在提交的数据中。
8.  **处理与通知**:
    *   工单按照定义的流程进行处理。
    *   当工单流程处理完成并到达结束节点时，后端系统会检查该工单是否关联了 `CreatorOpenId`。
    *   如果 OpenID 存在且微信通知已启用，系统将根据配置的 `template_id` 和 `template_format`，调用微信 API 向该 OpenID 发送处理结果的模板消息。

## 实现逻辑 (供维护参考)

*   **后端路由 (`router/router.go`)**:
    *   `GET /api/v1/wechat/auth`: 调用 `wechat.GetAuthUrl` 生成授权链接。
    *   `GET /api/v1/wechat/callback`: 调用 `wechat.Callback` 处理微信回调。
*   **微信服务 (`pkg/wechat/wechat.go`)**:
    *   `GetAuthUrl`: 根据配置生成微信网页授权 URL (`snsapi_base`)。
    *   `Callback`: 处理微信回调，用 `code` 换取 `openid`，设置 Cookie，重定向用户。
*   **微信通知服务 (`pkg/notify/wechat/wechat.go`)**:
    *   `GetAccessToken`: 从微信 API 获取 Access Token（**注意：当前无缓存**）。
    *   `SendWorkOrderResult`: 核心通知函数。检查配置，获取 Access Token，根据 `template_format` 构建消息体，调用微信 API 发送模板消息。
*   **前端工单创建 (`web/src/views/process/create-ticket.vue`)**:
    *   `created()` 和 `checkOpenid()`: 检查 Cookie 中的 `openid`，控制授权按钮显示。
    *   `handleWechatAuth()`: 点击授权按钮时调用后端获取授权 URL。
    *   `submitForm()`: 提交时从 Cookie 读取 `openid` 并附加到表单数据 (`formData.openid`)。
*   **后端工单创建 (`pkg/service/createWorkOrder.go`)**:
    *   在处理请求时，将前端传来的 `OpenId` 赋值给 `workOrderValue.CreatorOpenId`，随后保存到数据库。
*   **后端工单处理 (`pkg/service/handle.go`)**:
    *   在 `HandleWorkOrder` 方法中，当判断流程到达结束节点 (`h.targetStateValue["clazz"] == "end"`) 时，会检查 `h.workOrderDetails.CreatorOpenId` 是否为空，并异步 (`go func()`) 调用 `wechat.SendWorkOrderResult`。

## 注意事项与排错

*   **配置一致性**: 确保 `config/settings.yml` 中的配置（特别是 `appid`, `appsecret`, `redirect_uri`, `template_id`）与微信公众号后台的设置完全一致。
*   **域名问题**: `redirect_uri` 必须是已在微信后台配置的网页授权域名下的地址。
*   **IP 白名单**: 确保服务器 IP 已加入微信公众号后台的 IP 白名单。
*   **模板消息**: 确认模板消息功能已开通，模板 ID 正确，且模板内容与 `template_format` 配置的占位符匹配。
*   **日志**: 检查后端日志文件（具体路径取决于部署配置），特别是与微信相关的错误信息（如 Access Token 获取失败、发送模板消息失败等）。`pkg/notify/wechat/wechat.go` 中已添加相关日志输出。
*   **Access Token 缓存**: 当前实现未对 Access Token 进行缓存，频繁调用可能导致接口限制。如果遇到此问题，可以考虑引入 Redis 或内存缓存来存储和复用 Access Token。 