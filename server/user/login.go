package user

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go_crud/server/user/user_dao"
	"go_crud/server/user/utils"
	"go_crud/server/utils/token"
	"gorm.io/gorm"
	"time"
)

type LoginForm struct {
	Name     string
	Password string
}

func LoginPost(r *gin.RouterGroup, DB *gorm.DB) {
	r.POST("/login", func(c *gin.Context) {
		loginData := LoginForm{}
		err := c.ShouldBindJSON(&loginData)
		if err != nil { //数据错
			c.JSON(200, gin.H{
				"msg":  "数据校验未通过",
				"data": err.Error(),
				"code": "400",
			})
		} else {
			userDataList := user_dao.GetUserByName(loginData.Name, DB)
			if len(userDataList) == 0 { //没有查到
				c.JSON(200, gin.H{
					"msg":  "用户不存在",
					"data": loginData.Name,
					"code": "400",
				})
			} else {
				if time.Now().Before(userDataList[0].LockedUntil) {
					timeTemplate1 := "2006-01-02 15:04:05"
					c.JSON(200, gin.H{
						"msg":  "账户已被锁定到" + userDataList[0].LockedUntil.Format(timeTemplate1),
						"data": loginData.Name,
						"code": "400",
					})
				} else {
					rawPassword, _ := utils.RsaDecode(loginData.Password)
					loginPassword := utils.GetHash(rawPassword)
					if userDataList[0].Password == loginPassword {
						user_dao.RecordPasswordWrong(userDataList[0], DB, 0)
						user_dao.SetUserStatus(userDataList[0], DB, "in")
						// 签发token
						tokenDuration := time.Duration(viper.GetInt("token.shortDuration"))
						longDuration := time.Duration(viper.GetInt("token.longDuration"))
						signature, _ := token.IssueRS(loginData.Name, time.Now().Add(time.Minute*tokenDuration))
						c.Header("new_token", signature)
						signatureLong, _ := token.IssueRS(loginData.Name, time.Now().Add(time.Minute*longDuration))
						c.Header("new_long_token", signatureLong)
						c.Header("name", loginData.Name)
						c.JSON(200, gin.H{
							"msg":  "登录成功",
							"data": loginData.Name,
							"code": "233",
						})
					} else {
						user_dao.RecordPasswordWrong(userDataList[0], DB, userDataList[0].PasswordTry+1)
						c.JSON(200, gin.H{
							"msg":  "密码错误",
							"data": loginData.Name,
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
