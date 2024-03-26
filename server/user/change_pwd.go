package user

import (
	"github.com/gin-gonic/gin"
	"go_crud/server/user/user_dao"
	"go_crud/server/user/utils"
	"gorm.io/gorm"
	"time"
)

type ChangePwdForm struct {
	Name        string `json:"name"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func ChangePwdPost(r *gin.RouterGroup, DB *gorm.DB) {
	r.POST("/change_pwd", func(c *gin.Context) {
		Data := ChangePwdForm{}
		err := c.ShouldBindJSON(&Data)
		if err != nil { //数据错
			c.JSON(200, gin.H{
				"msg":  "数据校验未通过",
				"data": err.Error(),
				"code": "400",
			})
		} else {
			userDataList := user_dao.GetUserByName(Data.Name, DB)
			if len(userDataList) == 0 { //没有查到
				c.JSON(200, gin.H{
					"msg":  "用户不存在",
					"data": Data.Name,
					"code": "400",
				})
			} else {
				if time.Now().Before(userDataList[0].LockedUntil) {
					timeTemplate1 := "2006-01-02 15:04:05"
					c.JSON(200, gin.H{
						"msg":  "账户已被锁定到" + userDataList[0].LockedUntil.Format(timeTemplate1),
						"data": Data.Name,
						"code": "400",
					})
				} else {
					rawOldPassword, _ := utils.RsaDecode(Data.OldPassword)
					loginPassword := utils.GetHash(rawOldPassword)
					if userDataList[0].Password == loginPassword {
						rawNewPassword, _ := utils.RsaDecode(Data.NewPassword)
						if len(rawNewPassword) > 50 || len(rawNewPassword) < 8 {
							c.JSON(200, gin.H{
								"msg":  "密码长度需在8和50之间",
								"data": "",
								"code": "400",
							})
							return
						}
						newPassword := utils.GetHash(rawNewPassword)
						user_dao.SetUserPwd(userDataList[0], DB, newPassword)
						c.JSON(200, gin.H{
							"msg":  "修改成功",
							"data": Data.Name,
							"code": "234",
						})
					} else {
						user_dao.RecordPasswordWrong(userDataList[0], DB, userDataList[0].PasswordTry+1)
						c.JSON(200, gin.H{
							"msg":  "密码错误",
							"data": Data.Name,
							"code": "400",
						})
					}
				}
			}
		}
		//signature, _ := token.IssueHS("hello", time.Now().Add(time.Hour))
		//fmt.Println("签名内容", signature)
		//tokenErr := token.CheckHS(signature)
		//fmt.Println("验签", tokenErr == nil)

	})
}
