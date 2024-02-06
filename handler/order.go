package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/abdul-rehman-d/orders-api/model"
)

type Repo interface {
	Insert(ctx context.Context, order model.Order) error
	GetAll(ctx context.Context) ([]model.Order, error)
	Get(ctx context.Context, id int) (model.Order, error)
	Update(ctx context.Context, id int, order model.Order) (model.Order, error)
	Delete(ctx context.Context, id int) (model.Order, error)
}
type Order struct {
	Repo Repo
}

func (o *Order) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create order called")
}
func (o *Order) List(w http.ResponseWriter, r *http.Request) {
	fmt.Println("List order called")
}
func (o *Order) GetByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetByID order called")
}
func (o *Order) UpdateByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("UpdateByID order called")
}
func (o *Order) DeleteByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DeleteByID order called")
}
