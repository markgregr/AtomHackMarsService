package repository

import (
	"errors"
	"time"

	"github.com/SicParv1sMagna/AtomHackMarsService/internal/model"
)

func (r *Repository) GetDraftDocuments(page, pageSize int) ([]model.Document, error) {
    var documents []model.Document
    offset := (page - 1) * pageSize

    if err := r.db.DatabaseGORM.Where("status = ?", model.StatusDraft).Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&documents).Error; err != nil {
        return nil, err
    }

    return documents, nil
}

func (r *Repository) GetFormedDocuments(page, pageSize int) ([]model.Document, error) {
    var documents []model.Document
    offset := (page - 1) * pageSize

    if err := r.db.DatabaseGORM.Where("status = ?", model.StatusFormed).Order("sent_time DESC").Offset(offset).Limit(pageSize).Find(&documents).Error; err != nil {
        return nil, err
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
    if doc.DeliveryStatus != nil{
        existingDoc.DeliveryStatus = doc.DeliveryStatus
    }
    if doc.Status != ""{
        existingDoc.Status = doc.Status
    }
    if doc.SentTime != nil{
        existingDoc.SentTime = doc.SentTime
    }
    if doc.ReceivedTime != nil{
        existingDoc.ReceivedTime = doc.ReceivedTime
    }

    if err := r.db.DatabaseGORM.Save(existingDoc).Error; err != nil {
        return err
    }

    return nil
}

func (r *Repository) CreateDocument(doc *model.Document) (uint, error) {
	if err := r.db.DatabaseGORM.Create(doc).Error; err != nil {
		return 0, err
	}

	return doc.ID, nil
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

func (r *Repository) SendDocument(docID uint) (*model.Document, error) {
    doc, err := r.GetDocumentByID(docID)
    if err != nil {
        return nil, err
    }

    if doc.Status != model.StatusDraft {
        return nil, errors.New("document is not in draft status")
    }

    doc.Status = model.StatusFormed

    deliveryStatusPending := model.DeliveryStatusPending
    doc.DeliveryStatus = &deliveryStatusPending

    sentTime := time.Now()
    doc.SentTime = &sentTime

    if err := r.UpdateDocument(docID, doc); err != nil {
        return nil, err
    }

    return doc, nil
}

