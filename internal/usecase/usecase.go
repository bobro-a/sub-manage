package usecase

import (
	"context"
	"fmt"
	"log"
	"sub-manage/internal/repo"
	"sub-manage/pkg/models"
)

type SubUseCase interface {
	Create(ctx context.Context, sub models.Sub) (models.Sub, error)
	Read(ctx context.Context, id int64) (models.Sub, error)
	Update(ctx context.Context, sub models.Sub) (models.Sub, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]models.Sub, error)
	Sum(ctx context.Context, filter models.SumFilter) (int64, error)
}

type subUseCase struct {
	r repo.SubRepository
}

func New(r repo.SubRepository) SubUseCase {
	return &subUseCase{r: r}
}

func (uc *subUseCase) Create(ctx context.Context, sub models.Sub) (models.Sub, error) {
	if sub.ServiceName == "" || sub.Price <= 0 || sub.StartDate.After(sub.EndDate.Time) {
		return models.Sub{}, fmt.Errorf("incorrect request fields")
	}
	return uc.r.Create(ctx, sub)
}

func (uc *subUseCase) Read(ctx context.Context, id int64) (models.Sub, error) {
	log.Println("start usecase Read")
	if id <= 0 {
		return models.Sub{}, fmt.Errorf("invalid id")
	}
	log.Println("end usecase Read")
	return uc.r.Read(ctx, id)
}

func (uc *subUseCase) Update(ctx context.Context, sub models.Sub) (models.Sub, error) {
	if sub.ID <= 0 {
		return models.Sub{}, fmt.Errorf("invalid id")
	}
	if sub.ServiceName == "" || sub.Price <= 0 || sub.StartDate.After(sub.EndDate.Time) {
		return models.Sub{}, fmt.Errorf("incorrect request fields")
	}
	_, err := uc.Read(ctx, sub.ID)
	if err != nil {
		return models.Sub{}, fmt.Errorf("subscription not found")
	}
	return uc.r.Update(ctx, sub)
}

func (uc *subUseCase) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid id")
	}
	_, err := uc.Read(ctx, id)
	if err != nil {
		return fmt.Errorf("subscription not found")
	}
	return uc.r.Delete(ctx, id)
}

func (uc *subUseCase) List(ctx context.Context) ([]models.Sub, error) {
	return uc.r.List(ctx)
}

func (uc *subUseCase) Sum(ctx context.Context, filter models.SumFilter) (int64, error) {
	return uc.r.Sum(ctx, filter)
}
