package app

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

// Run запускает приложение
func (app *Application) Run() {
	r := gin.Default()

	ApiGroup := r.Group("/api/v1")
	{
		DocumentGroup := ApiGroup.Group("/document")
		{	
			DocumentGroup.POST("/", app.handler.CreateDocument)
			DocumentGroup.GET("/:docID", app.handler.GetDocumentByID)
			DocumentGroup.PUT("/:docID", app.handler.UpdateDocument)
			DocumentGroup.DELETE("/:docID", app.handler.DeleteDocument)
			DocumentGroup.POST("/:docID/file", app.handler.UploadFile)
			DocumentGroup.DELETE("/:docID/file/:fileID", app.handler.DeleteFile)

		}
	}

	WebSocketGroup := r.Group("/ws/v1")
	{
		DocumentGroup := WebSocketGroup.Group("/document")
		{
			DocumentGroup.GET("/", app.handler.GetDocuments)
		}

	}
	addr := fmt.Sprintf("%s:%d", app.cfg.API.ServiceHost, app.cfg.API.ServicePort)
	r.Run(addr)
	log.Println("Server down")
}
