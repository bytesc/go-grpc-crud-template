package crud_rpc

import (
	"context"
	"github.com/gin-gonic/gin"
	"go_crud/mysql_db"
	"go_crud/rpc_server/crud_pb"
	"gorm.io/gorm"
)

func QueryGET(r *gin.RouterGroup) {
	r.GET("/list/:name", func(c *gin.Context) {
		name := c.Param("name")
		var dataList []mysql_db.CrudList
		client := NewClient()
		result, err := client.Query(context.Background(), &crud_pb.QueryRequest{
			Name: name,
		})
		if err != nil {
			c.JSON(200, gin.H{
				"msg":  "查询失败",
				"data": err.Error(),
				"code": "400",
			})
		}

		if len(result.List) == 0 { //没有查到
			c.JSON(200, gin.H{
				"msg":  "查询失败，数据不存在",
				"data": dataList,
				"code": "400",
			})
		} else {
			for _, ql := range result.List {
				// 创建一个CrudList结构体实例
				cl := mysql_db.CrudList{
					// gorm.Model中的ID字段，这里假设ID是int类型
					Model: gorm.Model{
						ID: uint(ql.Id),
					},
					Name:     ql.Name,
					Level:    ql.Level,
					Email:    ql.Email,
					Phone:    ql.Phone,
					Birthday: ql.Birthday,
					Address:  ql.Address,
				}
				// 将CrudList结构体添加到切片中
				dataList = append(dataList, cl)
			}
			c.JSON(200, gin.H{
				"msg":  "查询成功",
				"data": dataList,
				"code": "200",
			})
		}
	})
}
