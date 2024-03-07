package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/SicParv1sMagna/AtomHackMarsService/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)


var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (h *Handler) GetDocuments(c *gin.Context) {
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade to websocket connection"})
        return
    }
    defer conn.Close()
    
    // Получаем параметры из URL
    page, _ := strconv.Atoi(c.Query("page"))
    pageSize, _ := strconv.Atoi(c.Query("pageSize"))
    statusStr  := c.Query("status")

	var status model.Status
    switch statusStr {
    case "Draft":
        status = model.StatusDraft
    case "Formed":
        status = model.StatusFormed
    default:
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
        return
    }

    for {
        // Получаем документы из репозитория
        documents, err := h.r.GetDocuments(page, pageSize, status)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve documents"})
            return
        }
        
        // Отправляем документы по WebSocket
        if err := conn.WriteJSON(documents); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send documents"})
            return
        }
    }
}

func (h *Handler) CreateDocument(c *gin.Context) {
	doc := model.Document{
		Status: model.StatusDraft,
		CreatedAt: time.Now(),
	}

	if err := h.r.CreateDocument(&doc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
        "message": "Document created successfully",
    })
}

func (h *Handler) GetDocumentByID(c *gin.Context) {
	docID, err := strconv.Atoi(c.Param("docID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file from request"})
		return
	}

	doc, err := h.r.GetDocumentByID(uint(docID)); 
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"doc": doc})
}

func (h *Handler) UpdateDocument(c *gin.Context) {
	docID, err := strconv.Atoi(c.Param("docID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file from request"})
		return
	}

    var doc model.Document

    if err := c.BindJSON(&doc); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind JSON"})
        return
    }

    if err := h.r.UpdateDocument(uint(docID), &doc); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "message": "Document updated successfully",
        "doc":     doc,
    })
}

func (h *Handler) DeleteDocument(c *gin.Context) {
    docID, err := strconv.Atoi(c.Param("docID"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get document ID from request"})
        return
    }

    if err := h.r.DeleteDocument(uint(docID)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete document: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Document deleted successfully"})
}




