package services

import (
	"errors"
	"net/http"
	"time"

	"em_subscription_manager/internal/apperrors"
	"em_subscription_manager/internal/models"
	"em_subscription_manager/internal/repo"

	"github.com/google/uuid"
)

type SubscriptionService interface {
	GetById(id int) (*models.Subscription, error)
	GetAll() ([]*models.Subscription, error)
	Create(subscription *models.Subscription) error
	Update(subscription *models.Subscription) error
	Delete(id int) error
	GetSumAmountByUserIdAndPeriodAndName(userId uuid.UUID, name string, startDate, endDate time.Time) (int, error)
}

type SubscriptionServiceImpl struct {
	repo repo.SubscriptionRepo
}

func NewSubscriptionService(repo repo.SubscriptionRepo) *SubscriptionServiceImpl {
	return &SubscriptionServiceImpl{repo: repo}
}

func (s *SubscriptionServiceImpl) GetById(id int) (*models.Subscription, error) {
	if id <= 0 {
		return nil, apperrors.NewBadRequest("id must be greater than zero")
	}

	return s.repo.GetById(id)
}

func (s *SubscriptionServiceImpl) GetAll() ([]*models.Subscription, error) {
	return s.repo.GetAll()
}

func (s *SubscriptionServiceImpl) Create(subscription *models.Subscription) error {
	if err := s.validate(subscription); err != nil {
		return err
	}

	if err := s.repo.Create(subscription); err != nil {
		return s.wrapRepoError(err, "failed to create subscription")
	}

	return nil
}

func (s *SubscriptionServiceImpl) Update(subscription *models.Subscription) error {
	if subscription == nil || subscription.Id <= 0 {
		return apperrors.NewBadRequest("subscription id must be provided")
	}

	if err := s.validate(subscription); err != nil {
		return err
	}

	if err := s.repo.Update(subscription); err != nil {
		return s.wrapRepoError(err, "failed to update subscription")
	}

	return nil
}

func (s *SubscriptionServiceImpl) Delete(id int) error {
	if id <= 0 {
		return apperrors.NewBadRequest("id must be greater than zero")
	}

	if err := s.repo.Delete(id); err != nil {
		return s.wrapRepoError(err, "failed to delete subscription")
	}

	return nil
}

func (s *SubscriptionServiceImpl) GetSumAmountByUserIdAndPeriodAndName(userId uuid.UUID, name string, startDate, endDate time.Time) (int, error) {
	if userId == uuid.Nil {
		return 0, apperrors.NewBadRequest("user_id is required")
	}
	if name == "" {
		return 0, apperrors.NewBadRequest("service_name is required")
	}
	if endDate.Before(startDate) {
		return 0, apperrors.NewBadRequest("end_date must be greater than or equal to start_date")
	}

	total, err := s.repo.GetSumAmountByUserIdAndPeriodAndName(userId, name, startDate, endDate)
	if err != nil {
		return 0, s.wrapRepoError(err, "failed to fetch subscription summary")
	}

	return total, nil
}

func (s *SubscriptionServiceImpl) validate(subscription *models.Subscription) error {
	if subscription == nil {
		return apperrors.NewBadRequest("subscription payload is required")
	}
	if subscription.ServiceName == "" {
		return apperrors.NewBadRequest("service_name is required")
	}
	if subscription.Price <= 0 {
		return apperrors.NewBadRequest("price must be greater than zero")
	}
	if subscription.UserId == uuid.Nil {
		return apperrors.NewBadRequest("user_id is required")
	}
	if subscription.StartDate.IsZero() {
		return apperrors.NewBadRequest("start_date is required")
	}
	if subscription.EndDate != nil && subscription.EndDate.Before(subscription.StartDate) {
		return apperrors.NewBadRequest("end_date must be after start_date")
	}

	return nil
}

func (s *SubscriptionServiceImpl) wrapRepoError(err error, message string) error {
	if err == nil {
		return nil
	}
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		return err
	}
	return apperrors.Wrap(err, http.StatusInternalServerError, message)
}
