package repo

import (
	"time"

	"em_subscription_manager/internal/models"

	"github.com/google/uuid"
)

type SubscriptionRepo interface {
	GetById(id int) (*models.Subscription, error)
	GetAll() ([]*models.Subscription, error)
	Create(subscription *models.Subscription) error
	Update(subscription *models.Subscription) error
	Delete(id int) error
	GetSumAmountByUserIdAndPeriodAndName(userId uuid.UUID, name string, startDate, endDate time.Time) (int, error)
}
