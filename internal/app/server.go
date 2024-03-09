package app

import (
	"fmt"
	"log"

	"github.com/SicParv1sMagna/AtomHackMarsService/docs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Run запускает приложение
func (app *Application) Run() {
	r := gin.Default()

	docs.SwaggerInfo.Title = "AtomHackMarsBackend RestAPI"
	docs.SwaggerInfo.Description = "API server for Mars application"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api/v1"

	// r.Use(middleware.CorsMiddleware())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	ApiGroup := r.Group("/api/v1")
	{
		DocumentGroup := ApiGroup.Group("/document")
		{
			DocumentGroup.GET("/draft", app.handler.GetDraftDocuments)
			DocumentGroup.GET("/formed", app.handler.GetFormedDocuments)
			DocumentGroup.POST("/", app.handler.CreateDocument)
			DocumentGroup.POST("/:docID", app.handler.SendDocument)
			DocumentGroup.GET("/:docID", app.handler.GetDocumentByID)
			DocumentGroup.PUT("/:docID", app.handler.UpdateDocument)
			DocumentGroup.DELETE("/:docID", app.handler.DeleteDocument)
			DocumentGroup.PUT("/:docID/file", app.handler.UploadFile)
			DocumentGroup.DELETE("/:docID/file/:fileID", app.handler.DeleteFile)
			DocumentGroup.PUT("/:docID/status", app.handler.UpdateStatusSuccess)

		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	addr := fmt.Sprintf("%s:%d", app.cfg.API.ServiceHost, app.cfg.API.ServicePort)
	r.Run(addr)

	log.Println("Server down")
}
