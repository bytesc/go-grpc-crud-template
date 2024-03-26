package crud_rpc

import (
	"context"
	"github.com/gin-gonic/gin"
	"go_crud/crud_rpc_server/crud_pb"
	"strconv"
)

func DeletePOST(r *gin.RouterGroup) {
	r.POST("/delete/:id", func(c *gin.Context) {

		idStr := c.Param("id") //接收路径参数
		// c.Query()  //接收查询参数
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(200, gin.H{
				"msg":  "id参数格式错误",
				"data": err.Error(),
				"code": "400",
			})
		}
		client := NewClient()
		result, err := client.Delete(context.Background(), &crud_pb.DeleteRequest{
			Id: int64(id),
		})
		if err != nil {
			c.JSON(200, gin.H{
				"msg":  "删除失败",
				"data": err.Error(),
				"code": "400",
			})
		}
		if result.Code != 0 {
			c.JSON(200, gin.H{
				"msg":  "删除失败",
				"data": result.Message,
				"code": "400",
			})
		}

		c.JSON(200, gin.H{
			"msg":  "删除成功",
			"data": idStr,
			"code": "200",
		})

	})
}
