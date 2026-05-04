package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"em_subscription_manager/internal/apperrors"
	"em_subscription_manager/internal/dto"
	"em_subscription_manager/internal/models"
	"em_subscription_manager/internal/services"
)

type SubscriptionHandler struct {
	service services.SubscriptionService
}

func RegisterSubscriptionRoutes(app *fiber.App, service services.SubscriptionService) {
	h := &SubscriptionHandler{service: service}

	app.Get("/subscriptions", h.GetAll)
	app.Get("/subscriptions/summary", h.GetSummary)
	app.Get("/subscriptions/:id", h.GetById)
	app.Post("/subscriptions", h.Create)
	app.Put("/subscriptions/:id", h.Update)
	app.Delete("/subscriptions/:id", h.Delete)
}

func respondWithError(c fiber.Ctx, err error) error {
	return c.Status(apperrors.GetCode(err)).JSON(dto.ErrorResponse{Error: apperrors.GetMessage(err)})
}

// GetAll godoc
// @Summary Get all subscriptions
// @Tags subscriptions
// @Produce json
// @Success 200 {array} models.Subscription
// @Failure 500 {object} dto.ErrorResponse
// @Router /subscriptions [get]
func (h *SubscriptionHandler) GetAll(c fiber.Ctx) error {
	subscriptions, err := h.service.GetAll()
	if err != nil {
		return respondWithError(c, err)
	}
	return c.JSON(subscriptions)
}

// GetById godoc
// @Summary Get subscription by id
// @Tags subscriptions
// @Produce json
// @Param id path int true "Subscription ID"
// @Success 200 {object} models.Subscription
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetById(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return respondWithError(c, apperrors.NewBadRequest("invalid id"))
	}

	subscription, err := h.service.GetById(id)
	if err != nil {
		return respondWithError(c, err)
	}
	return c.JSON(subscription)
}

// GetSummary godoc
// @Summary Get subscriptions summary
// @Tags subscriptions
// @Produce json
// @Param user_id query string true "User ID (UUID)"
// @Param service_name query string false "Service name"
// @Param start_date query string true "Start date (MM-YYYY)"
// @Param end_date query string true "End date (MM-YYYY)"
// @Success 200 {object} map[string]int
// @Router /subscriptions/summary [get]
func (h *SubscriptionHandler) GetSummary(c fiber.Ctx) error {
	userID, err := uuid.Parse(c.Query("user_id"))
	if err != nil {
		return respondWithError(c, apperrors.NewBadRequest("invalid user_id"))
	}

	service_name := c.Query("service_name")

	startDate, err := models.ParseMonthYear(c.Query("start_date"))
	if err != nil {
		return respondWithError(c, apperrors.NewBadRequest("invalid start_date; expected MM-YYYY"))
	}

	endDate, err := models.ParseMonthYear(c.Query("end_date"))
	if err != nil {
		return respondWithError(c, apperrors.NewBadRequest("invalid end_date; expected MM-YYYY"))
	}

	total, err := h.service.GetSumAmountByUserIdAndPeriodAndName(userID, service_name, startDate, endDate)
	if err != nil {
		return respondWithError(c, err)
	}

	return c.JSON(fiber.Map{"total": total})
}

// Create godoc
// @Summary Create subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param request body dto.CreateSubscriptionDTO true "Subscription payload"
// @Success 201 {object} models.Subscription
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /subscriptions [post]
func (h *SubscriptionHandler) Create(c fiber.Ctx) error {
	var dtoBody dto.CreateSubscriptionDTO
	if err := c.Bind().JSON(&dtoBody); err != nil {
		return respondWithError(c, apperrors.NewBadRequest("invalid request body or date format (MM-YYYY)"))
	}

	subscription := &models.Subscription{
		ServiceName: dtoBody.ServiceName,
		Price:       dtoBody.Price,
		UserId:      dtoBody.UserId,
		StartDate:   dtoBody.StartDate,
		EndDate:     dtoBody.EndDate,
	}

	if err := h.service.Create(subscription); err != nil {
		return respondWithError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(subscription)
}

// Update godoc
// @Summary Update subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path int true "Subscription ID"
// @Param request body dto.UpdateSubscriptionDTO true "Subscription payload"
// @Success 200 {object} models.Subscription
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) Update(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return respondWithError(c, apperrors.NewBadRequest("invalid id"))
	}

	var dtoBody dto.UpdateSubscriptionDTO
	if err := c.Bind().JSON(&dtoBody); err != nil {
		return respondWithError(c, apperrors.NewBadRequest("invalid request body"))
	}

	subscription := &models.Subscription{
		Id:          id,
		ServiceName: dtoBody.ServiceName,
		Price:       dtoBody.Price,
		UserId:      dtoBody.UserId,
		StartDate:   dtoBody.StartDate,
		EndDate:     dtoBody.EndDate,
	}

	if err := h.service.Update(subscription); err != nil {
		return respondWithError(c, err)
	}

	return c.JSON(subscription)
}

// Delete godoc
// @Summary Delete subscription
// @Tags subscriptions
// @Param id path int true "Subscription ID"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) Delete(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return respondWithError(c, apperrors.NewBadRequest("invalid id"))
	}

	if err := h.service.Delete(id); err != nil {
		return respondWithError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
