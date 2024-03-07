package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UploadFile обрабатывает запрос на загрузку файла.
// @Summary Загружает файл.
// @Description Загружает файл в хранилище MinIO и связывает его с указанным документом.
// @Tags Файлы
// @Accept multipart/form-data
// @Produce json
// @Param docID path int true "Идентификатор документа"
// @Param file formData file true "Файл для загрузки"
// @Success 200 {object} model.MessageResponse "Успешный ответ"
// @Failure 400 {object} model.ErrorResponse "Ошибка в запросе"
// @Failure 500 {object} model.ErrorResponse "Внутренняя ошибка сервера"
// @Router /document/{docID}/file [post]
func (h *Handler) UploadFile(c *gin.Context) {
	docID, err := strconv.Atoi(c.Param("docID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file from request"})
		return
	}

	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file from request"})
		return
	}
	defer file.Close()

	// Получаем размер файла
	fileSize := fileHeader.Size
	fileName := fileHeader.Filename

	// Загружаем файл в хранилище MinIO и обновляем путь файла в базе данных
	if err := h.r.UploadFile(uint(docID), file, fileSize, fileName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})
}

// DeleteFile обрабатывает запрос на удаление файла.
// @Summary Удаляет файл.
// @Description Удаляет файл из хранилища MinIO и из базы данных.
// @Tags Файлы
// @Accept json
// @Produce json
// @Param docID path int true "Идентификатор документа"
// @Param fileID path int true "Идентификатор файла"
// @Success 200 {object} model.MessageResponse "Успешный ответ"
// @Failure 400 {object} model.ErrorResponse "Ошибка в запросе"
// @Failure 500 {object} model.ErrorResponse "Внутренняя ошибка сервера"
// @Router /document/{docID}/file/{fileID} [delete]
func (h *Handler) DeleteFile(c *gin.Context) {
	docID, err := strconv.Atoi(c.Param("docID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get document ID from request"})
		return
	}

	fileID, err := strconv.Atoi(c.Param("fileID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file ID from request"})
		return
	}

	if err := h.r.DeleteFileByID(uint(docID), uint(fileID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}
