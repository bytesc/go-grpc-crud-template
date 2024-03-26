package crud_rpc

import (
	"context"
	"github.com/gin-gonic/gin"
	"go_crud/mysql_db"
	"go_crud/rpc_server/crud_pb"
	"net/url"
	"strconv"
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

		if result.Code != 0 {
			c.JSON(200, gin.H{
				"msg":  "查询失败",
				"data": result.Message,
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
			mysql_db.CrudListRpcToOrm(result.List, &dataList)
			c.JSON(200, gin.H{
				"msg":  "查询成功",
				"data": dataList,
				"code": "200",
			})
		}
	})
}

// QueryPageGET 分页查询
func QueryPageGET(r *gin.RouterGroup) {
	r.GET("/list/", func(c *gin.Context) {
		var dataList []mysql_db.CrudList
		var pageSize, pageNum int
		pageSizeStr := c.Query("pageSize")
		pageNumStr := c.Query("pageNum")
		if pageSizeStr != "" {
			pageSizeInt, err := strconv.Atoi(pageSizeStr) //?pageSize=xxx 查询参数
			pageSize = pageSizeInt
			// 127.0.0.1:8080/crud/list/?pageNum=2&pageSize=3
			if err != nil {
				c.JSON(200, gin.H{
					"msg":  "查询失败，pageNum参数格式错误",
					"data": err.Error(),
					"code": "400",
				})
			}
		} else {
			pageSize = -1
		}
		if pageNumStr != "" {
			pageNumInt, err := strconv.Atoi(pageNumStr)
			pageNum = pageNumInt
			if err != nil {
				c.JSON(200, gin.H{
					"msg":  "查询失败，pageSize参数格式错误",
					"data": err.Error(),
					"code": "400",
				})
			}
		} else {
			pageNum = 1
		}

		client := NewClient()
		result, err := client.QueryPage(context.Background(), &crud_pb.QueryPageRequest{
			PageNum:     int32(pageNum),
			PageSize:    int32(pageSize),
			QueryParams: valuesToStringMap(c.Request.URL.Query()),
		})
		if err != nil {
			c.JSON(200, gin.H{
				"msg":  "查询失败",
				"data": err.Error(),
				"code": "400",
			})
		}

		if result.Code != 0 {
			c.JSON(200, gin.H{
				"msg":  "查询失败",
				"data": result.Message,
				"code": "400",
			})
		}

		mysql_db.CrudListRpcToOrm(result.List, &dataList)

		if len(dataList) == 0 { //没有查到
			c.JSON(200, gin.H{
				"msg":  "查询失败，数据不存在",
				"data": dataList,
				"code": "400",
			})
		} else {
			c.JSON(200, gin.H{
				"msg": "查询成功",
				"data": gin.H{
					"list":     dataList,
					"total":    result.Total,
					"pageNum":  pageNum,
					"pageSize": pageSize,
				},
				"code": "200",
			})
		}
	})
}

func valuesToStringMap(values url.Values) map[string]string {
	stringMap := make(map[string]string)
	for key, values := range values {
		if key != "pageSize" && key != "pageNum" {
			if len(values) > 0 {
				stringMap[key] = values[0]
			}
		}
	}
	return stringMap
}