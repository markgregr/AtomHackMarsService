package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/SicParv1sMagna/AtomHackMarsService/internal/model"
	"github.com/gin-gonic/gin"
)

// GetDraftDocuments возвращает черновики документов.
// @Summary Возвращает черновики документов.
// @Description Возвращает список черновиков документов с учетом параметров page и pageSize.
// @Tags Документы
// @Accept json
// @Produce json
// @Param page query int false "Номер страницы" default(1)
// @Param pageSize query int false "Размер страницы" default(10)
// @Success 200 {array} model.GetDocuments "Успешный ответ"
// @Failure 500 {object} model.ErrorResponse "Внутренняя ошибка сервера"
// @Router /document/draft [get]
func (h *Handler) GetDraftDocuments(c *gin.Context) {
    page, _ := strconv.Atoi(c.Query("page"))
    pageSize, _ := strconv.Atoi(c.Query("pageSize"))

    // Получаем черновики документов
    documents, total, err := h.r.GetDraftDocuments(page, pageSize)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve documents: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"items":documents,
								"total": total})
}

// GetFormedDocuments возвращает сформированные документы.
// @Summary Возвращает сформированные документы.
// @Description Возвращает список сформированных документов с учетом параметров page и pageSize.
// @Tags Документы
// @Accept json
// @Produce json
// @Param page query int false "Номер страницы" default(1)
// @Param pageSize query int false "Размер страницы" default(10)
// @Success 200 {array} model.GetDocuments "Успешный ответ"
// @Failure 500 {object} model.ErrorResponse "Внутренняя ошибка сервера"
// @Router /document/formed [get]
func (h *Handler) GetFormedDocuments(c *gin.Context) {
    // Получаем параметры из URL
    page, _ := strconv.Atoi(c.Query("page"))
    pageSize, _ := strconv.Atoi(c.Query("pageSize"))

    // Получаем сформированные документы
    documents, total, err := h.r.GetFormedDocuments(page, pageSize)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve documents: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"items":documents,
								"total": total})
}

// CreateDocument создает новый документ.
// @Summary Создает новый документ.
// @Description Создает новый документ на основе переданных данных JSON, возвращает id созданного документа.
// @Tags Документы
// @Accept json
// @Produce json
// @Success 200 {object} model.DocumentCreate "Успешный ответ"
// @Failure 400 {object} model.ErrorResponse "Ошибка в запросе"
// @Failure 500 {object} model.ErrorResponse "Внутренняя ошибка сервера"
// @Router /document [post]
func (h *Handler) CreateDocument(c *gin.Context) {
	doc := model.Document{
		Status:    model.StatusDraft,
		CreatedAt: time.Now(),
	}

	docID, err := h.r.CreateDocument(&doc)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id": docID,
	})
}

// GetDocumentByID получает документ по его ID.
// @Summary Получает документ по ID.
// @Description Получает документ из репозитория по указанному ID.
// @Tags Документы
// @Accept json
// @Produce json
// @Param docID path int true "ID документа"
// @Success 200 {object} model.Document "Успешный ответ"
// @Failure 400 {object} model.ErrorResponse "Ошибка в запросе"
// @Failure 500 {object} model.ErrorResponse "Внутренняя ошибка сервера"
// @Router /document/{docID} [get]
func (h *Handler) GetDocumentByID(c *gin.Context) {
	docID, err := strconv.Atoi(c.Param("docID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file from request: " + err.Error()})
		return
	}

	doc, err := h.r.GetDocumentByID(uint(docID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, doc)
}

// UpdateDocument обновляет информацию о документе.
// @Summary Обновляет информацию о документе.
// @Description Обновляет информацию о документе на основе переданных данных JSON.
// @Tags Документы
// @Accept json
// @Produce json
// @Param docID path int true "ID документа"
// @Param doc body model.DocumentUpdate true "Пользовательский объект в формате JSON"
// @Success 200 {object} model.MessageResponse "Успешный ответ"
// @Failure 400 {object} model.ErrorResponse "Ошибка в запросе"
// @Failure 500 {object} model.ErrorResponse "Внутренняя ошибка сервера"
// @Router /document/{docID} [put]
func (h *Handler) UpdateDocument(c *gin.Context) {
	docID, err := strconv.Atoi(c.Param("docID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file from request: " + err.Error()})
		return
	}

	var doc model.Document

	if err := c.BindJSON(&doc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind JSON: " + err.Error()})
		return
	}

	if err := h.r.UpdateDocument(uint(docID), &doc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document updated successfully"})
}

// DeleteDocument удаляет документ по его ID.
// @Summary Удаляет документ по ID.
// @Description Удаляет документ из репозитория по указанному ID.
// @Tags Документы
// @Accept json
// @Produce json
// @Param docID path int true "ID документа"
// @Success 200 {object} model.MessageResponse "Успешный ответ"
// @Failure 400 {object} model.ErrorResponse "Ошибка в запросе"
// @Failure 500 {object} model.ErrorResponse "Внутренняя ошибка сервера"
// @Router /document/{docID} [delete]
func (h *Handler) DeleteDocument(c *gin.Context) {
	docID, err := strconv.Atoi(c.Param("docID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get document ID from request: " + err.Error()})
		return
	}

	if err := h.r.DeleteDocument(uint(docID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete document: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document deleted successfully"})
}

// SendDocument отправляет документ на Землю.
// @Summary Отправляет документ на Землю.
// @Description Отправляет документ на Землю по docID.
// @Tags Документы
// @Accept json
// @Produce json
// @Param docID path int true "ID документа"
// @Success 200 {object} model.MessageResponse "Успешный ответ"
// @Failure 400 {object} model.ErrorResponse "Ошибка в запросе"
// @Failure 500 {object} model.ErrorResponse "Внутренняя ошибка сервера"
// @Router /document/{docID} [post]
func (h *Handler) SendDocument(c *gin.Context) {
	docID, err := strconv.Atoi(c.Param("docID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get document ID from request: " + err.Error()})
		return
	}

	doc, err := h.r.GetDocumentByID(uint(docID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get document from database: " + err.Error()})
		return
	}

	docSubmitted, err := h.r.SendDocument(uint(docID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send document: " + err.Error()})
		return
	}

	docJSON, err := json.Marshal(docSubmitted)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize document: " + err.Error()})
		return
	}

	// Отправка документа в Kafka
	if err = h.p.SendReport(h.p.KafkaCfg.Topic, string(docJSON)); err != nil {
		if err := h.r.UpdateDocument(uint(docID), doc); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update document message: " + err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document sent to Kafka successfully"})
}
