package crud_rpc

import (
	"context"
	"github.com/gin-gonic/gin"
	"go_crud/crud_rpc_server/crud_pb"
	"go_crud/mysql_db"
	"strconv"
)

func UpdatePOST(r *gin.RouterGroup) {
	r.POST("/update/:id/", func(c *gin.Context) {
		idStr := c.Param("id") //接收路径参数
		// c.Query()  //接收查询参数
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(200, gin.H{
				"msg":  "id格式错误",
				"data": err,
				"code": "400",
			})
		}
		if err != nil {
			c.JSON(200, gin.H{
				"msg":  "id参数格式错误",
				"data": err.Error(),
				"code": "400",
			})
		}

		var listRes mysql_db.CrudList
		err = c.ShouldBindJSON(&listRes) //数据校验
		if err != nil {                  //数据错
			c.JSON(200, gin.H{
				"msg":  "添加失败，数据校验未通过",
				"data": err.Error(),
				"code": "400",
			})
			//fmt.Println(err)
		}
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
		result, err := client.Update(context.Background(), &crud_pb.UpdateRequest{
			Id:   int64(id),
			List: pbCL,
		})
		if err != nil {
			c.JSON(200, gin.H{
				"msg":  "更新失败",
				"data": err.Error(),
				"code": "400",
			})
		}
		if result.Code != 0 {
			c.JSON(200, gin.H{
				"msg":  "更新失败",
				"data": result.Message,
				"code": "400",
			})
		} else {
			c.JSON(200, gin.H{
				"msg":  "更新成功",
				"data": idStr,
				"code": "200",
			})
		}
	})
}
