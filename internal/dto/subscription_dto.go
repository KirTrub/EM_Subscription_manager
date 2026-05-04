package dto

import (
	"em_subscription_manager/internal/models"

	"github.com/google/uuid"
)

type CreateSubscriptionDTO struct {
	ServiceName string                   `json:"service_name"`
	Price       int                      `json:"price"`
	UserId      uuid.UUID                `json:"user_id"`
	StartDate   models.SubscriptionDate  `json:"start_date"`
	EndDate     *models.SubscriptionDate `json:"end_date,omitempty"`
}

type UpdateSubscriptionDTO struct {
	ServiceName string                   `json:"service_name"`
	Price       int                      `json:"price"`
	UserId      uuid.UUID                `json:"user_id"`
	StartDate   models.SubscriptionDate  `json:"start_date"`
	EndDate     *models.SubscriptionDate `json:"end_date,omitempty"`
}
