package routers

import (
	"github.com/letscrum/letscrum/middlewares/jwt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"github.com/letscrum/letscrum/pkg/export"
	"github.com/letscrum/letscrum/pkg/qrcode"
	"github.com/letscrum/letscrum/pkg/upload"
	"github.com/letscrum/letscrum/routers/apis"
	"github.com/letscrum/letscrum/routers/apis/v1"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.StaticFS("/export", http.Dir(export.GetExcelFullPath()))
	r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))
	r.StaticFS("/qrcode", http.Dir(qrcode.GetQrCodeFullPath()))

	r.POST("/auth", apis.GetAuth)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/upload", apis.UploadImage)

	apisV1 := r.Group("/apis/v1")
	apisV1.Use()
	{
		apisV1.GET("/signin", v1.SignIn)
	}
	apisV1.Use(jwtMiddleware.JWT())
	{
		//获取标签列表
		apisV1.GET("/tags", v1.GetTags)
		//新建标签
		apisV1.POST("/tags", v1.AddTag)
		//更新指定标签
		apisV1.PUT("/tags/:id", v1.EditTag)
		//删除指定标签
		apisV1.DELETE("/tags/:id", v1.DeleteTag)
		//导出标签
		r.POST("/tags/export", v1.ExportTag)
		//导入标签
		r.POST("/tags/import", v1.ImportTag)

		//获取文章列表
		apisV1.GET("/articles", v1.GetArticles)
		//获取指定文章
		apisV1.GET("/articles/:id", v1.GetArticle)
		//新建文章
		apisV1.POST("/articles", v1.AddArticle)
		//更新指定文章
		apisV1.PUT("/articles/:id", v1.EditArticle)
		//删除指定文章
		apisV1.DELETE("/articles/:id", v1.DeleteArticle)
		//生成文章海报
		apisV1.POST("/articles/poster/generate", v1.GenerateArticlePoster)
		apisV1.POST("/projects", v1.CreateProject)
		apisV1.GET("/projects", v1.ListProject)
		apisV1.PUT("/projects/:name", v1.UpdateProject)
		apisV1.DELETE("/projects/:name", v1.DeleteProject)
		apisV1.GET("/projects/:name", v1.GetProject)
	}

	return r
}
