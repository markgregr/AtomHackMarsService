package repository

import (
	"fmt"
	"mime/multipart"

	"github.com/SicParv1sMagna/AtomHackMarsService/internal/model"
)

func (r *Repository) GetFilesByDocumentID(docID uint) ([]model.File, error) {
	var files []model.File
	if err := r.db.DatabaseGORM.Where("document_id = ?", docID).Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}

func (r *Repository) UploadFile(docID uint, file multipart.File, fileSize int64, fileName string) error {
	var document model.Document
	if err := r.db.DatabaseGORM.First(&document, docID).Error; err != nil {
		return fmt.Errorf("failed to find document with ID %d: %w", docID, err)
	}

	objectName := fmt.Sprintf("documents/%d/%s", docID, fileName)

	fileURL, err := r.mc.UploadFile(objectName, file, fileSize)
	if err != nil {
		return err
	}

	newFile := model.File{
		Path:       fileURL,
		DocumentID: docID,
	}
	// Сохраняем новую запись файла в базе данных
	if err := r.db.DatabaseGORM.Create(&newFile).Error; err != nil {
		return fmt.Errorf("failed to save file path in database: %w", err)
	}

	return nil
}

func (r *Repository) DeleteFileByID(docID, fileID uint) error {
    var file model.File
    if err := r.db.DatabaseGORM.Where("document_id = ? AND id = ?", docID, fileID).First(&file).Error; err != nil {
        return fmt.Errorf("failed to find file with ID %d and document ID %d: %w", fileID, docID, err)
    }

    if err := r.mc.DeleteFile(file.Path); err != nil {
        return err
    }

    if err := r.db.DatabaseGORM.Delete(&file).Error; err != nil {
        return fmt.Errorf("failed to delete file from database: %w", err)
    }

    return nil
}

func (r *Repository) DeleteFilesByDocumentID(docID uint) error {
	files, err := r.GetFilesByDocumentID(docID)
	if err != nil {
		return err
	}

	for _, file := range files {
		if err := r.DeleteFileByID(docID, file.ID); err != nil {
			return err
		}
	}

	return nil
}
