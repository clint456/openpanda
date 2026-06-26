// ============================================================
// 包名: utils
// 说明: 统一返回格式封装
//
//	所有接口都使用此格式返回，前端只需处理这一种数据结构
//	拓展方式：如需新增返回字段（如分页信息），修改 Response 结构体即可
//
// ============================================================
package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一API返回结构
// 前端接收到的 JSON 格式：
//
//	{
//	  "code": 200,         // 业务状态码
//	  "message": "success", // 提示信息
//	  "data": {...}         // 实际数据（可为 nil）
//	}
type Response struct {
	Code    int         `json:"code"`    // 业务状态码，200=成功
	Message string      `json:"message"` // 提示信息
	Data    interface{} `json:"data"`    // 数据载荷，interface{} 表示任意类型
}

// Success 成功响应（带数据）
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 成功响应（自定义消息）
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: message,
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, httpStatus int, code int, message string) {
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// BadRequest 400 错误（参数错误）
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, 400, message)
}

// Unauthorized 401 错误（未授权）
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, 401, message)
}

// NotFound 404 错误
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, 404, message)
}

// InternalError 500 服务器内部错误
func InternalError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, 500, message)
}

// ============================================================
// 分页返回结构（预留）
// 后续列表接口需要分页时，使用此结构体
// ============================================================

// PageData 分页数据结构
type PageData struct {
	List     interface{} `json:"list"`      // 数据列表
	Total    int64       `json:"total"`     // 总记录数
	Page     int         `json:"page"`      // 当前页码
	PageSize int         `json:"page_size"` // 每页条数
}

// SuccessWithPage 成功响应（带分页）
func SuccessWithPage(c *gin.Context, data PageData) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}
