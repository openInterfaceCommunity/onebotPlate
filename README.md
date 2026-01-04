<div align="center">

<img src="https://p.ipic.vip/mml1wv.jpg" alt="OneBot Plate Logo" width="400" />

<div style="font-size: 2.5em; font-weight: bold; margin: 20px 0;">OneBot Plate (Platform)</div>

### **[onebotPlate](https://github.com/openInterfaceCommunity/onebotPlate)** 是一个基于 [OneBot 12 标准](https://github.com/botuniverse/onebot/tree/main/specs) 开发的机器人平台实现。

</div>

## 核心概念

在使用本项目前，请确保您已了解 OneBot 12 中的以下核心概念：

*   **机器人平台 (Bot Platform / platform)**：提供聊天机器人 API 的社交软件平台。
    *   例如：`qq`、`discord`、`kook`、`telegram`、`wechat` 等。
*   **OneBot 实现 (OneBot Implementation / impl)**：负责对接具体的机器人平台，并向上层（应用端）提供符合 OneBot 12 标准接口的程序。**OneBot Plate 即属于此类实现。**
*   **OneBot 应用 (OneBot Application)**：与实现端交互，处理机器人具体业务逻辑的程序（例如使用 NoneBot2 编写的插件）。
*   **机器人 (Bot)**：指代具体的机器人账号。在 OneBot 12 中，由 `self_id` 和 `platform` 共同唯一确定。

## 快速入门

### 1. 注册与登录
访问项目托管的 Web 控制台（Dashboard）进行账户操作：
*   **控制台地址**: [OneBot Plate Dashboard](https://incitymega.cn/bot/devtool/dashboard)
*   **测试账户**:
    *   Alice (密码: `pass123`)
    *   Bob (密码: `pass123`)

### 2. 创建 Bot 实例
1.  登录控制台后，进入 **BotManager** 页面。
2.  点击注册（创建）新的 Bot。
3.  获取该 Bot 的 **Access Token**，这是后续连接的必要凭据。

### 3. 开发第一个机器机器人
项目提供了 Python 和 Go 的示例代码，您可以参考这些例子快速开始：

*   **WebSocket 模式 (Echo Bot)**: [examples/echo_bot.py](file:///examples/echo_bot.py)
*   **HTTP 模式 (Minimal Bot)**: [examples/http_bot.py](file:///examples/http_bot.py)

> [!TIP]
> 仅 **Bot** 身份支持 WebSocket 事件推送，以保证用户与机器人的权责划分清晰。Bot 需要由用户账号在控制台中创建。

## 开发者指南

### API 交互方式
*   **Dashboard API Console**: 您可以在控制台的 API 调试板块中实时测试所有 API。
*   **认证方式**: 连接时需通过 `access_token` 进行鉴权。
*   **通信模式**:
    *   **HTTP**: 适用于主动调用 Action。
    *   **WebSocket**: 适用于监听 Event 以及接收 Action 回调。

### 特别说明
*   **独立实体**: 用户和机器人是平台中的独立实体，但机器人必须由用户创建。
*   **多实例支持**: 同一个认证信息支持多个链接实例。
*   **群组支持**: 机器人可以被邀请进入群组，并实时监听群组消息。

## 社区与支持

如果您有任何问题或建议，欢迎通过以下方式联系：

**项目主页**: [GitHub Repo](https://github.com/openInterfaceCommunity/onebotPlate)

**QQ 交流群**:1077671218

<img src="https://p.ipic.vip/93dvt2.jpg" alt="QQ Group QR Code" width="200" />


### 项目支持

如果您觉得本项目对您有所帮助，欢迎进行打赏或赞助支持：

<img src="https://p.ipic.vip/t2s20z.jpg" alt="Donation QR Code" width="200" />

### 项目费用

由于服务器资源有限:收取少许合理费用  
账号(不删档):暂定100元  
(注:目前测试无需费用,但可能因为人多,导致服务性能下降而被删档)   

个人咨询QQ:306582825
商业咨询微信:extent