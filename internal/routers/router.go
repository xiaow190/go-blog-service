package routers

import (
	_ "go-programming-book/blog-service/docs"
	"go-programming-book/blog-service/global"
	"go-programming-book/blog-service/internal/middleware"
	api "go-programming-book/blog-service/internal/routers/api"
	v1 "go-programming-book/blog-service/internal/routers/api/v1"
	"go-programming-book/blog-service/pkg/limiter"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

var methodLimiters = limiter.NewMethodLimiter().AddBuckets(
	limiter.LimiterBucketRule{
		Key:          "/auth",
		FillInterval: time.Second,
		Capacity:     10,
		Quantum:      10,
	},
)

func NewRouter() *gin.Engine {
	r := gin.New()
	if global.ServerSetting.RunMode == "debug" {
		r.Use(gin.Logger())
		r.Use(gin.Recovery())
	} else {
		r.Use(middleware.AccessLog())
		r.Use(middleware.Recovery())
	}
	r.Use(middleware.Tracering())
	r.Use(middleware.RateLimiter(methodLimiters))
	r.Use(middleware.ContextTimeout(global.AppSetting.DefaultContextTimeout))
	r.Use(middleware.Translations())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	article := v1.NewArticle()
	tag := v1.NewTag()
	upload := NewUpload()
	// 文件上传
	r.POST("/upload/file", upload.uploadFile)
	r.StaticFS("/static", http.Dir(global.AppSetting.UploadSavePath))
	//  jwt
	r.POST("/auth", api.GetAuth)
	apiv1 := r.Group("/api/v1")
	apiv1.Use(middleware.JWT())
	{
		// 标签相关接口
		apiv1.POST("/tags", tag.Create)
		apiv1.DELETE("/tags/:id", tag.Delete)
		apiv1.PUT("/tags/:id", tag.Update)
		apiv1.PATCH("/tags/:id/state", tag.Update)
		apiv1.GET("/tags", tag.List)
		apiv1.GET("/tags/:id", tag.Get)

		// 文章
		apiv1.POST("/articles", article.Create)
		apiv1.DELETE("/articles/:id", article.Delete)
		apiv1.PUT("/articles/:id", article.Update)
		apiv1.PATCH("/articles/:id/state", article.Update)
		apiv1.GET("/articles/:id", article.Get)
		apiv1.GET("/articles", article.List)
	}
	return r
}
