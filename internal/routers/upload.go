package routers

import (
	"go-programming-book/blog-service/global"
	"go-programming-book/blog-service/internal/service"
	"go-programming-book/blog-service/pkg/app"
	"go-programming-book/blog-service/pkg/convert"
	"go-programming-book/blog-service/pkg/errcode"
	"go-programming-book/blog-service/pkg/upload"

	"github.com/gin-gonic/gin"
)

type Upload struct{}

func NewUpload() Upload {
	return Upload{}
}

func (u Upload) uploadFile(c *gin.Context) {
	response := app.NewResponse(c)
	file, fileHeader, err := c.Request.FormFile("file")
	fileType := convert.StrTo(c.PostForm("type")).MustInt()
	if err != nil {
		errRsp := errcode.InvalidParams.WithDetails(err.Error())
		response.ToErrorResponse(errRsp)
		return
	}
	if fileHeader == nil || fileType <= 0 {
		response.ToErrorResponse(errcode.InvalidParams)
		return
	}
	svc := service.New(c.Request.Context())
	fileInfo, err := svc.UploadFile(upload.FileType(fileType), file, fileHeader)
	if err != nil {
		global.Logger.Errorf(c, "svc.UploadFile err: %v", err)
		errRsp := errcode.ErrorUploadFileFaill.WithDetails(err.Error())
		response.ToErrorResponse(errRsp)
		return
	}

	response.ToResponse(gin.H{"file_access_url": fileInfo.AccessUrl})
	return
}
