package repository

import (
	"errors"

	"github.com/SicParv1sMagna/AtomHackMarsService/internal/model"
)

func (r *Repository) GetDocuments(page, pageSize int, status model.Status) ([]model.Document, error) {
    var documents []model.Document
    offset := (page - 1) * pageSize

    switch status {
    case model.StatusDraft:
        // Для документов со статусом "Черновик" сортируем по дате создания
        if err := r.db.DatabaseGORM.Where("status = ?", model.StatusDraft).Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&documents).Error; err != nil {
            return nil, err
        }
    case model.StatusFormed:
        // Для документов со статусом "Сформирован" сортируем по SentTime
        if err := r.db.DatabaseGORM.Where("status = ?", model.StatusFormed).Order("sent_time DESC").Offset(offset).Limit(pageSize).Find(&documents).Error; err != nil {
            return nil, err
        }
    default:
        return nil, errors.New("invalid status")
    }

    return documents, nil
}

func (r *Repository) GetDocumentByID(docID uint) (*model.Document, error) {
    var document model.Document
    if err := r.db.DatabaseGORM.Preload("Files").Where("status != ?", model.StatusDeleted).First(&document, docID).Error; err != nil {
        return nil, err
    }
    return &document, nil
}

func (r *Repository) UpdateDocument(docID uint, doc *model.Document) error {
    existingDoc, err := r.GetDocumentByID(docID)
    if err != nil {
        return err
    }

    if doc.Title != "" {
        existingDoc.Title = doc.Title
    }
    if doc.Owner != "" {
        existingDoc.Owner = doc.Owner
    }
    if doc.Payload != "" {
        existingDoc.Payload = doc.Payload
    }

    if err := r.db.DatabaseGORM.Save(existingDoc).Error; err != nil {
        return err
    }

    return nil
}

func (r *Repository) CreateDocument(doc *model.Document) error {
	if err := r.db.DatabaseGORM.Create(doc).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteDocument(docID uint) error {
    doc, err := r.GetDocumentByID(docID)
    if err != nil {
        return err
    }

    doc.Status = model.StatusDeleted

    if err := r.UpdateDocument(docID, doc); err != nil {
        return err
    }

    return nil
}

