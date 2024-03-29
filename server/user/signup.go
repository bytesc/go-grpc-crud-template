package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go_crud/mysql_db"
	"go_crud/server/user/user_dao"
	"go_crud/server/user/utils"
	"time"
)

type SignupForm struct {
	Name     string `validate:"required,alphanum,min=3,max=50"`
	Email    string `validate:"required,email,max=50"`
	Password string `validate:"required"`
}

func SignUpPost(r *gin.RouterGroup) {
	r.POST("/signup", func(c *gin.Context) {
		signupData := SignupForm{}
		err := c.ShouldBindJSON(&signupData)
		if err != nil { //数据错
			c.JSON(200, gin.H{
				"msg":  "数据校验未通过",
				"data": err.Error(),
				"code": "400",
			})
		} else {
			validate := validator.New()
			err := validate.Struct(signupData)
			if err != nil {
				c.JSON(200, gin.H{
					"msg":  "数据验证不通过，用户名只能包含数字和字母，邮箱需有效",
					"data": err.Error(),
					"code": "400",
				})
			} else {
				userDataList := user_dao.GetUserByName(signupData.Name)
				if len(userDataList) != 0 { //查到
					c.JSON(200, gin.H{
						"msg":  "用户已经存在",
						"data": signupData.Name,
						"code": "400",
					})
				} else {
					var userData mysql_db.UserList
					userData.Name = signupData.Name
					userData.Email = signupData.Email
					userData.LockedUntil = time.Now()
					rawPassword, _ := utils.RsaDecode(signupData.Password)
					fmt.Println(rawPassword, len(rawPassword))
					if len(rawPassword) > 50 || len(rawPassword) < 8 {
						c.JSON(200, gin.H{
							"msg":  "密码长度需在8和50之间",
							"data": "",
							"code": "400",
						})
						return
					}
					userData.Password = utils.GetHash(rawPassword)
					//hashBytes := sha256.Sum256([]byte(signupData.Password))
					result := user_dao.CreateUser(userData)
					if result.Error != nil {
						c.JSON(200, gin.H{
							"msg":  "注册失败",
							"data": result.Error,
							"code": "400",
						})
					} else {
						c.JSON(200, gin.H{
							"msg":  "注册成功",
							"data": userData.Name,
							"code": "234",
						})
					}
				}
			}
		}

	})
}
