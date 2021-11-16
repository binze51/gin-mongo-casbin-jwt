package response

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RetCode int

const (
	//操作成功的状态码
	Success RetCode = 0

	//系统功能错误
	ServerError         RetCode = 1001
	TooManyRequests     RetCode = 1002
	ParamBindError      RetCode = 1003
	AuthorizationError  RetCode = 1004
	CallHTTPError       RetCode = 1005
	ParamUploadError    RetCode = 1006
	UploadFileSizeError RetCode = 1007
	FileReadError       RetCode = 1008
	FileOpenError       RetCode = 1009
	FileTypeError       RetCode = 1010
	FileSaveError       RetCode = 1011
	LoadFileError       RetCode = 1012
	TokenExpiredError   RetCode = 1013
)

var CodeTable = map[RetCode]string{
	Success: "操作成功",

	ServerError:         "服务器内部错误",
	TooManyRequests:     "请求数过多",
	ParamBindError:      "参数提交错误",
	AuthorizationError:  "授权认证失败",
	CallHTTPError:       "调用第三方 HTTP 接口失败",
	ParamUploadError:    "头像上传参数错误",
	UploadFileSizeError: "头像文件大于5m错误",
	FileReadError:       "头像文件读流失败",
	FileOpenError:       "头像文件打开失败",
	FileTypeError:       "头像文件仅支持png或jpg类型",
	FileSaveError:       "头像文件保存失败",
	LoadFileError:       "获取文件失败",
	TokenExpiredError:   "授权token已过期",
}

func MergeCodeTable(othermodel map[RetCode]string) {
	for k, v := range othermodel {
		if CodeTable[k] == "" {
			CodeTable[k] = v
			continue
		}
		log.Panicf("othermodel[%d]:[%s] 重复", k, v)
		return
	}
}

type ResqData struct {
	BizCode RetCode     `json:"bizeCode"`       // 业务操作返回码，为0时 没有错误
	Message string      `json:"msg"`            // 描述信息
	Data    interface{} `json:"data,omitempty"` // 数据
}

func Result(BizCode RetCode, msg string, data interface{}, c *gin.Context) {
	// 使用bizCode来判断，任何时候都返回200
	c.JSON(http.StatusOK, ResqData{
		BizCode,
		msg,
		data,
	})
}

func Ok(c *gin.Context) {
	Result(Success, CodeTable[Success], nil, c)
}
func OkWithData(data interface{}, c *gin.Context) {
	Result(Success, CodeTable[Success], data, c)
}

func FileLoad(filebyte []byte, c *gin.Context) {
	c.Writer.Write(filebyte)
}
func Fail(errcode RetCode, c *gin.Context) {
	Result(errcode, CodeTable[errcode], nil, c)
}
func FailWithErrorRaw(errcode RetCode, message string, c *gin.Context) {
	Result(errcode, CodeTable[errcode]+":"+message, nil, c)
}
func FailWithError(errcode RetCode, message string, c *gin.Context) {
	Result(errcode, message, nil, c)
}

func CodeToError(code RetCode) error {
	if CodeTable[code] == "" {
		panic("请完善你新增的错误情况")
	}
	return errors.New(CodeTable[code])
}
