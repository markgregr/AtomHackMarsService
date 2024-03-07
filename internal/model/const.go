package model

// DeliveryStatus представляет статус доставки
type DeliveryStatus string

const (
	DeliveryStatusSuccess DeliveryStatus = "SUCCESS"
	DeliveryStatusPending DeliveryStatus = "PENDING"
	DeliveryStatusError   DeliveryStatus = "ERROR"
)

// Status представляет статус документа
type Status string

const (
	StatusDeleted Status = "DELETED"
	StatusDraft   Status = "DRAFT"
	StatusFormed  Status = "FORMED"
)