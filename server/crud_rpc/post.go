package crud_rpc

import (
	"context"
	"github.com/gin-gonic/gin"
	"go_crud/mysql_db"
	"go_crud/rpc_server/crud_pb"
)

func AddPOST(r *gin.RouterGroup) {
	r.POST("/add", func(c *gin.Context) {

		var listRes mysql_db.CrudList
		err := c.ShouldBindJSON(&listRes) //数据校验
		if err != nil {                   //数据错
			c.JSON(200, gin.H{
				"msg":  "添加失败，数据校验未通过",
				"data": err.Error(),
				"code": "400",
			})
			//fmt.Println(err)
		} else {
			pbCL := &crud_pb.CrudList{
				Id:       int64(listRes.ID),
				Name:     listRes.Name,
				Level:    listRes.Level,
				Email:    listRes.Email,
				Phone:    listRes.Phone,
				Birthday: listRes.Birthday,
				Address:  listRes.Address,
			}

			client := NewClient()
			result, err := client.Add(context.Background(), &crud_pb.AddRequest{
				List: pbCL,
			})
			if err != nil {
				c.JSON(200, gin.H{
					"msg":  "添加失败",
					"data": err.Error(),
					"code": "400",
				})
			}
			if result.Code != 0 {
				c.JSON(200, gin.H{
					"msg":  "添加失败",
					"data": result.Message,
					"code": "400",
				})
			} else {
				c.JSON(200, gin.H{
					"msg":  "添加成功",
					"data": "",
					"code": "200",
				})
			}

		}
	})
	return
}
