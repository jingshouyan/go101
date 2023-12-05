package util

import (
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

const defaultPageSize = 10

func GetLimitOffset(c *gin.Context) (limit, offset int) {
	page, _ := com.StrTo(c.Query("page")).Int()
	limit, _ = com.StrTo(c.Query("pageSize")).Int()
	if limit <= 0 {
		limit = defaultPageSize
	}
	if page > 0 {
		offset = (page - 1) * limit
	}
	return
}

func PageResult(total int64, list interface{}) map[string]interface{} {
	page := make(map[string]interface{})
	page["total"] = total
	page["list"] = list
	return page
}
