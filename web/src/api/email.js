import service from '@/utils/request'

// @Tags email
// @Summary 发送测试邮件
// @Security ApiKeyAuth
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"返回成功"}"
// @Router /email/emailTest [post]
export const emailTest = (data) => {
    return service({
        url: "/system/emailTest",
        method: 'post',
        data
    })
}

export const botTest = () => {
    return service({
        url: "/system/botTest",
        method: "get"
    })
}