package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sub-manage/internal/usecase"
	"sub-manage/pkg/models"
	"time"

	"github.com/google/uuid"
)

const (
	timeForRequest = 2
	parseType      = "2006-01"
)

type subHandler struct {
	uc usecase.SubUseCase
}

func NewSubHandler(uc usecase.SubUseCase) *subHandler {
	return &subHandler{uc: uc}
}

// @Summary      Create a new subscription
// @Description  Creates a new subscription with the provided data
// @Accept       json
// @Produce      json
// @Param        subscription body models.Sub true "Subscription object to create"
// @Success      200 {object} models.Sub
// @Failure      400 {string} string "Invalid request fields"
// @Router       / [post]
func (sh *subHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), timeForRequest*time.Second)
	defer cancel()

	var sub models.Sub
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	created, err := sh.uc.Create(ctx, sub)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

// @Summary      Get a subscription of subscriptions
// @Description  Gets a single subscription by ID
// @Accept       json
// @Produce      json
// @Param        id query int false "ID of the subscription"
// @Success      200 {object} models.Sub
// @Failure      400 {string} string "Invalid ID format"
// @Failure      404 {string} string "Subscription not found"
// @Router       / [get]
func (sh *subHandler) Read(w http.ResponseWriter, r *http.Request) {
	log.Println("start handler Read")
	ctx, cancel := context.WithTimeout(r.Context(), timeForRequest*time.Second)
	defer cancel()

	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	sub, err := sh.uc.Read(ctx, id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sub)
	log.Println("end handler Read")
}

// @Summary      Update an existing subscription
// @Description  Updates a subscription with the provided data. The ID in the body is required.
// @Accept       json
// @Produce      json
// @Param        subscription body models.Sub true "Subscription object with updated values"
// @Success      200 {object} models.Sub
// @Failure      400 {string} string "Invalid ID or request fields"
// @Failure      404 {string} string "Subscription not found"
// @Router       / [put]
func (sh *subHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), timeForRequest*time.Second)
	defer cancel()

	var sub models.Sub
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	updated, err := sh.uc.Update(ctx, sub)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

// @Summary      Delete a subscription
// @Description  Deletes a subscription by ID.
// @Accept       json
// @Produce      json
// @Param        id query int true "ID of the subscription to delete"
// @Success      204 {string} string "No Content"
// @Failure      400 {string} string "Invalid ID format"
// @Failure      404 {string} string "Subscription not found"
// @Router       / [delete]
func (sh *subHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), timeForRequest*time.Second)
	defer cancel()

	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = sh.uc.Delete(ctx, id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// @Summary      Get a list of all subscriptions
// @Description  Retrieves all subscriptions from the database.
// @Accept       json
// @Produce      json
// @Success      200 {array} models.Sub
// @Router       /list [get]
func (sh *subHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), timeForRequest*time.Second)
	defer cancel()

	subs, err := sh.uc.List(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subs)
}

// @Summary      Calculate the total cost
// @Description  Calculates the total cost of subscriptions based on user_id, service_name, start_date, and end_date.
// @Accept       json
// @Produce      json
// @Param        user_id query string false "User ID to filter subscriptions" format(uuid)
// @Param        service_name query string false "Service name to filter subscriptions"
// @Param        start_date query string false "Start date for calculation (YYYY-MM)"
// @Param        end_date query string false "End date for calculation (YYYY-MM)"
// @Success      200 {object} object "Returns a JSON object with the total cost"
// @Failure      400 {string} string "Invalid query parameters"
// @Router       /sum [get]
func (sh *subHandler) Sum(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), timeForRequest*time.Second)
	defer cancel()

	var filter models.SumFilter
	fromStr := r.URL.Query().Get("start_date")
	toStr := r.URL.Query().Get("end_date")
	if fromStr != "" {
		from, err := time.Parse(parseType, fromStr)
		if err != nil {
			http.Error(w, "invalid from date", http.StatusBadRequest)
			return
		}
		filter.StartDate = &models.MonthYear{Time: from}
	}
	if toStr != "" {
		to, err := time.Parse(parseType, toStr)
		if err != nil {
			http.Error(w, "invalid to date", http.StatusBadRequest)
			return
		}
		filter.EndDate = &models.MonthYear{Time: to}
	}

	if userIDStr := r.URL.Query().Get("user_id"); userIDStr != "" {
		parsedUUID, err := uuid.Parse(userIDStr)
		if err != nil {
			http.Error(w, "invalid user_id format", http.StatusBadRequest)
			return
		}
		filter.UserID = &parsedUUID
	}

	if nameStr := r.URL.Query().Get("service_name"); nameStr != "" {
		filter.ServiceName = nameStr
	}

	total, err := sh.uc.Sum(ctx, filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int64{"total_cost": total})
}
