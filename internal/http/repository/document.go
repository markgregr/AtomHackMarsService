package repository

import (
	"errors"
	"strings"
	"time"

	"github.com/SicParv1sMagna/AtomHackMarsService/internal/model"
)


func (r *Repository) GetDocumentsFormedCount(status model.Status, deliveryStatus model.DeliveryStatus, ownerOrTitle string) (uint, error) {
    var count int64
    query := r.db.DatabaseGORM.Model(&model.Document{})

    if status != "" {
        query = query.Where("LOWER(status) = LOWER(?)", status)
    }

    if deliveryStatus != "" {
        query = query.Where("LOWER(delivery_status) = LOWER(?)", deliveryStatus)
    }

    if ownerOrTitle != "" {
        query = query.Where("owner LIKE ? OR LOWER(title) LIKE ?", "%"+ownerOrTitle+"%", "%"+strings.ToLower(ownerOrTitle)+"%")
    }

    if err := query.Count(&count).Error; err != nil {
        return 0, err
    }

    return uint(count), nil
}

func (r *Repository) GetDocumentsDraftCount(status model.Status, deliveryStatus model.DeliveryStatus, owner, title string) (uint, error) {
    var count int64
    query := r.db.DatabaseGORM.Model(&model.Document{})

    if status != "" {
        query = query.Where("LOWER(status) = LOWER(?)", status)
    }

    if deliveryStatus != "" {
        query = query.Where("LOWER(delivery_status) = LOWER(?)", deliveryStatus)
    }

    if owner != "" {
        query = query.Where("owner =", owner)
    }

    if title != "" {
        query = query.Where("LOWER(title) LIKE ?", "%"+strings.ToLower(title)+"%")
    }

    if err := query.Count(&count).Error; err != nil {
        return 0, err
    }

    return uint(count), nil
}

func (r *Repository) GetFormedDocuments(page, pageSize int, deliveryStatus model.DeliveryStatus, ownerOrTitle string) ([]model.Document, uint, error) {
    var documents []model.Document
    offset := (page - 1) * pageSize

    query := r.db.DatabaseGORM

    if deliveryStatus != "" {
        query = query.Where("LOWER(delivery_status) = LOWER(?)", deliveryStatus)

        if deliveryStatus == model.DeliveryStatusSuccess {
            query = query.Order("received_time DESC")
        } else if deliveryStatus == model.DeliveryStatusPending {
            query = query.Order("sent_time DESC")
        }
    } else {
        query = query.Where("LOWER(status) = LOWER(?)", model.StatusFormed).Order("sent_time DESC")
    }

    if ownerOrTitle != "" {
        query = query.Where("owner LIKE ? OR LOWER(title) LIKE ?", "%"+ownerOrTitle+"%", "%"+strings.ToLower(ownerOrTitle)+"%")
    }

    if err := query.Offset(offset).Limit(pageSize).Find(&documents).Error; err != nil {
        return nil, 0, err
    }

    total, err := r.GetDocumentsFormedCount(model.StatusFormed, deliveryStatus, ownerOrTitle)
    if err != nil {
        return nil, 0, err
    }

    return documents, total, nil
}

func (r *Repository) GetDraftDocuments(page, pageSize int, owner, title string) ([]model.Document, uint, error) {
    var documents []model.Document
    offset := (page - 1) * pageSize

    query := r.db.DatabaseGORM.Where("LOWER(status) = LOWER(?)", model.StatusDraft)

    if owner != "" {
        query = query.Where("owner =", owner)
    }

    if title != "" {
        title = strings.ToLower(title)
        query = query.Where("LOWER(title) LIKE ?", "%"+title+"%")
    }

    if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&documents).Error; err != nil {
        return nil, 0, err
    }

    total, err := r.GetDocumentsDraftCount(model.StatusDraft, "", owner, title)
    if err != nil {
        return nil, 0, err
    }

    return documents, total, nil
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

