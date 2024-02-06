package handler

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/abdul-rehman-d/orders-api/model"
	"github.com/abdul-rehman-d/orders-api/repository/order"
	"github.com/google/uuid"
)

type Repo interface {
	Insert(ctx context.Context, order model.Order) error
	GetAll(ctx context.Context, page order.Page) (order.GetAllResult, error)
	Get(ctx context.Context, id uint64) (model.Order, error)
	Update(ctx context.Context, id uint64, order model.Order) error
	Delete(ctx context.Context, id uint64) error
}
type Order struct {
	Repo Repo
}

var ErrNotFound = errors.New("order does not exist")

func formatMessage(key, str string) []byte {
	message := map[string]string{
		key: str,
	}
	out, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	return out
}

func formatSuccessMessage(str string) []byte {
	return formatMessage("message", str)
}

func formatErrorMessage(str string) []byte {
	return formatMessage("error", str)
}

func returnEarlyWithMessage(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	w.Write(formatErrorMessage(message))
}

type CreateReq struct {
	CustomerID uuid.UUID        `json:"customer_id"`
	LineItems  []model.LineItem `json:"line_items"`
}

func validateBody(body CreateReq) bool {
	if body.CustomerID.ID() == 0 {
		return false
	}
	if len(body.LineItems) == 0 {
		return false
	}
	for _, item := range body.LineItems {
		if item.ItemID.ID() == 0 {
			return false
		}
		if item.Quantity <= 0 {
			return false
		}
	}
	return true
}

func (o *Order) Create(w http.ResponseWriter, r *http.Request) {
	var body CreateReq
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil || !validateBody(body) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(formatErrorMessage("invalid body"))
		return
	}

	now := time.Now().UTC()

	order := model.Order{
		OrderID:    rand.Uint64(),
		CustomerID: body.CustomerID,
		LineItems:  body.LineItems,
		CreatedAt:  &now,
	}

	if err = o.Repo.Insert(r.Context(), order); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(formatErrorMessage("something went wrong"))
		return
	}

	res, err := json.Marshal(order)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(formatErrorMessage("something went wrong"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

func (o *Order) List(w http.ResponseWriter, r *http.Request) {
	cursorStr := r.URL.Query().Get("page")
	if cursorStr == "" {
		cursorStr = "0"
	}

	cursor, err := strconv.ParseUint(cursorStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(formatErrorMessage("page can only be int"))
		return
	}

	result, err := o.Repo.GetAll(r.Context(), order.Page{
		Cursor: cursor,
		Count:  10,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(formatErrorMessage("something went wrong"))
		return
	}

	response, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(formatErrorMessage("something went wrong"))
		return
	}

	w.Write(response)
}

func (o *Order) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(formatErrorMessage("id can only be int"))
		return
	}

	order, err := o.Repo.Get(r.Context(), id)

	if err != nil {
		if errors.Is(err, ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.Write(formatErrorMessage(err.Error()))
		return
	}

	res, err := json.Marshal(order)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(formatErrorMessage("something went wrong"))
		return
	}

	w.Write(res)
}

func (o *Order) UpdateByID(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(formatErrorMessage("invalid body"))
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(formatErrorMessage("id can only be int"))
		return
	}

	currentOrder, err := o.Repo.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.Write(formatErrorMessage(err.Error()))
		return
	}

	now := time.Now().UTC()

	switch body.Status {
	case "shipped":
		if currentOrder.ShippedAt != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(formatErrorMessage("invalid body"))
			return
		}
		currentOrder.ShippedAt = &now
	case "completed":
		if currentOrder.ShippedAt == nil || currentOrder.CompletedAt != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(formatErrorMessage("invalid body"))
			return
		}
		currentOrder.CompletedAt = &now
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write(formatErrorMessage("invalid body"))
		return
	}

	err = o.Repo.Update(r.Context(), id, currentOrder)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(formatErrorMessage("something went wrong"))
		return
	}

	res, err := json.Marshal(currentOrder)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(formatErrorMessage("something went wrong"))
		return
	}

	w.Write(res)
}

func (o *Order) DeleteByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(formatErrorMessage("id can only be int"))
		return
	}

	err = o.Repo.Delete(r.Context(), id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(formatErrorMessage(err.Error()))
		return
	}

	w.Write(formatSuccessMessage("order deleted"))
}
