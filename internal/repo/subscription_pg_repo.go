package repo

import (
	"database/sql"
	"net/http"
	"time"

	"em_subscription_manager/internal/apperrors"
	"em_subscription_manager/internal/models"

	"github.com/google/uuid"
)

type SubscriptionPgRepo struct {
	db *sql.DB
}

func New(db *sql.DB) *SubscriptionPgRepo {
	return &SubscriptionPgRepo{db: db}
}

func (r *SubscriptionPgRepo) GetById(id int) (*models.Subscription, error) {
	var subscription models.Subscription
	var endDate sql.NullTime

	err := r.db.QueryRow(
		"SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions WHERE id = $1",
		id,
	).Scan(&subscription.Id, &subscription.ServiceName, &subscription.Price, &subscription.UserId, &subscription.StartDate, &endDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperrors.NewNotFound("subscription not found")
		}
		return nil, apperrors.Wrap(err, http.StatusInternalServerError, "failed to fetch subscription")
	}
	if endDate.Valid {
		sd := models.SubscriptionDate(endDate.Time)
		subscription.EndDate = &sd
	}

	return &subscription, nil
}

func (r *SubscriptionPgRepo) GetAll() ([]*models.Subscription, error) {
	rows, err := r.db.Query("SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions")
	if err != nil {
		return nil, apperrors.Wrap(err, http.StatusInternalServerError, "failed to fetch subscriptions")
	}
	defer rows.Close()

	var subscriptions []*models.Subscription
	for rows.Next() {
		var subscription models.Subscription
		var endDate sql.NullTime

		err := rows.Scan(&subscription.Id, &subscription.ServiceName, &subscription.Price, &subscription.UserId, &subscription.StartDate, &endDate)
		if err != nil {
			return nil, apperrors.Wrap(err, http.StatusInternalServerError, "failed to read subscription row")
		}
		if endDate.Valid {
			sd := models.SubscriptionDate(endDate.Time)
			subscription.EndDate = &sd
		}
		subscriptions = append(subscriptions, &subscription)
	}

	return subscriptions, nil
}

func (r *SubscriptionPgRepo) Create(subscription *models.Subscription) error {
	_, err := r.db.Exec(
		"INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date) VALUES ($1, $2, $3, $4, $5)",
		subscription.ServiceName,
		subscription.Price,
		subscription.UserId,
		subscription.StartDate,
		subscription.EndDate,
	)
	if err != nil {
		return apperrors.Wrap(err, http.StatusInternalServerError, "failed to create subscription")
	}
	return nil
}

func (r *SubscriptionPgRepo) Update(subscription *models.Subscription) error {
	result, err := r.db.Exec(
		"UPDATE subscriptions SET service_name = $1, price = $2, user_id = $3, start_date = $4, end_date = $5 WHERE id = $6",
		subscription.ServiceName,
		subscription.Price,
		subscription.UserId,
		subscription.StartDate,
		subscription.EndDate,
		subscription.Id,
	)
	if err != nil {
		return apperrors.Wrap(err, http.StatusInternalServerError, "failed to update subscription")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return apperrors.Wrap(err, http.StatusInternalServerError, "failed to update subscription")
	}
	if rowsAffected == 0 {
		return apperrors.NewNotFound("subscription not found")
	}

	return nil
}

func (r *SubscriptionPgRepo) Delete(id int) error {
	result, err := r.db.Exec("DELETE FROM subscriptions WHERE id = $1", id)
	if err != nil {
		return apperrors.Wrap(err, http.StatusInternalServerError, "failed to delete subscription")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return apperrors.Wrap(err, http.StatusInternalServerError, "failed to delete subscription")
	}
	if rowsAffected == 0 {
		return apperrors.NewNotFound("subscription not found")
	}

	return nil
}

func (r *SubscriptionPgRepo) GetSumAmountByUserIdAndPeriodAndName(userId uuid.UUID, serviceName string, startDate, endDate time.Time) (int, error) {
	query := `
		SELECT COALESCE(SUM(sub_total), 0) FROM (
			SELECT 
				price * (
					EXTRACT(YEAR FROM age(LEAST(COALESCE(end_date, $2), $2), GREATEST(start_date, $3))) * 12 +
					EXTRACT(MONTH FROM age(LEAST(COALESCE(end_date, $2), $2), GREATEST(start_date, $3))) + 1
				) as sub_total
			FROM subscriptions 
			WHERE user_id = $1 
			  AND start_date <= $2 
			  AND (end_date IS NULL OR end_date >= $3)
			  AND ($4 = '' OR LOWER(service_name) = LOWER($4))
		) as totals`

	var total float64
	err := r.db.QueryRow(query, userId, endDate, startDate, serviceName).Scan(&total)
	if err != nil {
		return 0, apperrors.Wrap(err, http.StatusInternalServerError, "failed to fetch summary")
	}

	return int(total), nil
}
