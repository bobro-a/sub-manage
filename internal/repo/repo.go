package repo

import (
	"context"
	"log"
	"sub-manage/pkg/models"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type SubRepository interface {
	Create(ctx context.Context, sub models.Sub) (models.Sub, error)
	Read(ctx context.Context, id int64) (models.Sub, error)
	Update(ctx context.Context, sub models.Sub) (models.Sub, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]models.Sub, error)
	Sum(ctx context.Context, filter models.SumFilter) (int64, error)
}

type subRepo struct {
	db        *sqlx.DB
	tableName string
	psql      squirrel.StatementBuilderType
}

func NewSubRepo(tableName string, db *sqlx.DB) SubRepository {
	return &subRepo{tableName: tableName,
		db:   db,
		psql: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)}
}

func (s *subRepo) Create(ctx context.Context, sub models.Sub) (models.Sub, error) {
	query := s.psql.Insert(s.tableName).
		Columns("service_name", "price", "user_id", "start_date", "end_date").
		Values(sub.ServiceName, sub.Price, sub.UserID, sub.StartDate.Time, sub.EndDate.Time).
		Suffix("RETURNING id")

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return sub, err
	}

	err = s.db.QueryRowContext(ctx, sqlStr, args...).Scan(&sub.ID)
	return sub, err
}

func (s *subRepo) Read(ctx context.Context, id int64) (models.Sub, error) {
	log.Println("start Repo Read")
	var sub models.Sub
	query := s.psql.Select("id", "service_name", "price", "user_id", "start_date", "end_date").
		From(s.tableName).
		Where(squirrel.Eq{"id": id})

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return sub, err
	}
	err = s.db.GetContext(ctx, &sub, sqlStr, args...)
	log.Println("end Repo Read")
	return sub, err
}

func (s *subRepo) Update(ctx context.Context, sub models.Sub) (models.Sub, error) {
	query := s.psql.Update(s.tableName).
		Set("service_name", sub.ServiceName).
		Set("price", sub.Price).
		Set("user_id", sub.UserID).
		Set("start_date", sub.StartDate.Time).
		Set("end_date", sub.EndDate.Time).
		Where(squirrel.Eq{"id": sub.ID})
	sqlStr, args, err := query.ToSql()
	if err != nil {
		return sub, err
	}
	_, err = s.db.ExecContext(ctx, sqlStr, args...)
	return sub, err
}

func (s *subRepo) Delete(ctx context.Context, id int64) error {
	query := s.psql.Delete(s.tableName).
		Where(squirrel.Eq{"id": id})

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx, sqlStr, args...)
	return err
}

func (s *subRepo) List(ctx context.Context) ([]models.Sub, error) {
	log.Println("start Repo List")

	var sub []models.Sub
	query := s.psql.Select("*").From(s.tableName)
	sqlStr, args, err := query.ToSql()
	if err != nil {
		return sub, err
	}

	err = s.db.SelectContext(ctx, &sub, sqlStr, args...)
	log.Println("end Repo List")
	return sub, err
}

func (s *subRepo) Sum(ctx context.Context, filter models.SumFilter) (int64, error) {
	query := s.psql.Select("COALESCE(SUM(price),0)").From(s.tableName)

	if filter.StartDate != nil { // Проверяем, что дата была задана
		query = query.Where(squirrel.GtOrEq{"start_date": filter.StartDate.Time})
	}

	if filter.EndDate != nil { // Проверяем, что дата была задана
		query = query.Where(squirrel.LtOrEq{"end_date": filter.EndDate.Time})
	}

	if filter.UserID != nil {
		query = query.Where(squirrel.Eq{"user_id": *filter.UserID})
	}

	if filter.ServiceName != "" {
		query = query.Where(squirrel.Eq{"service_name": filter.ServiceName})
	}
	sqlStr, args, err := query.ToSql()
	if err != nil {
		return 0, err
	}
	var total int64
	err = s.db.QueryRowContext(ctx, sqlStr, args...).Scan(&total)
	return total, err
}
