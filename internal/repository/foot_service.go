package repository

import (
	"carWash/internal/domain"
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"reflect"
	"strings"
	"time"
)

type FootServiceRepos struct {
	db *sqlx.DB
}

func (f *FootServiceRepos) Create(ctx *fiber.Ctx, service domain.FootService) (int, error) {
	var id int

	_, cancel := context.WithTimeout(ctx.Context(), 4*time.Second)

	defer cancel()

	query := fmt.Sprintf(
		`INSERT 
					INTO 
					%s 
						(service_name, price)
					VALUES
						($1,$2)
			 RETURNING id`, serviceTable)

	err := f.db.QueryRowx(query, service.ServiceName, service.Price).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("repository.Create: %w", err)
	}
	return id, nil
}

func (f *FootServiceRepos) GetAll(ctx *fiber.Ctx, page domain.Pagination) (*domain.GetAllResponses, error) {

	_, cancel := context.WithTimeout(ctx.Context(), 4*time.Second)

	defer cancel()

	count, err := countPage(f.db, serviceTable, "")
	if err != nil {
		return nil, fmt.Errorf("repository.GetAll: %w", err)
	}

	offset, pagesCount := calculatePagination(&page, count)

	inp := make([]*domain.FootService, 0, page.Limit)

	query := fmt.Sprintf("SELECT * FROM %s ORDER BY id ASC LIMIT $1 OFFSET $2", serviceTable)

	err = f.db.Select(&inp, query, page.Limit, offset)

	if err != nil {
		return nil, fmt.Errorf("repository.GetAll: %w", err)
	}

	pages := domain.PaginationPage{
		Page:  page.Page,
		Pages: pagesCount,
		Count: count,
	}

	ans := domain.GetAllResponses{
		Data:     inp,
		PageInfo: pages,
	}
	return &ans, nil
}

func (f *FootServiceRepos) GetById(ctx *fiber.Ctx, id int) (*domain.FootService, error) {
	var (
		place domain.FootService
	)

	_, cancel := context.WithTimeout(ctx.Context(), 4*time.Second)

	defer cancel()

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", serviceTable)

	err := f.db.Get(&place, query, id)

	if err != nil {
		return nil, fmt.Errorf("repository.GetById: %w", domain.ErrNotFound)
	}

	return &place, nil
}

func (f *FootServiceRepos) Update(ctx *fiber.Ctx, id int, inp domain.FootService) error {
	_, cancel := context.WithTimeout(ctx.Context(), 4*time.Second)

	defer cancel()

	setValues := make([]string, 0, reflect.TypeOf(domain.FootService{}).NumField())

	if inp.Price != 0 {
		setValues = append(setValues, fmt.Sprintf("price=:price"))
	}

	if inp.ServiceName != "" {
		setValues = append(setValues, fmt.Sprintf("service_name=:service_name"))
	}

	setQuery := strings.Join(setValues, ", ")

	if setQuery == "" {
		return fmt.Errorf("repository.Update: %w", errors.New("empty body"))
	}

	query := fmt.Sprintf(`UPDATE %s SET %s WHERE id = %d`, serviceTable, setQuery, id)

	result, err := f.db.NamedExec(query, inp)

	if err != nil {
		return fmt.Errorf("repository.Update: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository.Update: %w", err)
	}

	if affected == 0 {
		return fmt.Errorf("repository.Update: %w", domain.ErrNotFound)
	}

	return nil
}

func (f *FootServiceRepos) Delete(ctx *fiber.Ctx, id int) error {

	_, cancel := context.WithTimeout(ctx.Context(), 500*time.Millisecond)

	defer cancel()

	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1 ", serviceTable)

	err := f.db.QueryRowx(query, id).Err()

	if err != nil {
		return fmt.Errorf("repository.Delete: %w", domain.ErrNotFound)
	}
	return nil

}

func NewFootServiceRepos(db *sqlx.DB) *FootServiceRepos {
	return &FootServiceRepos{db: db}
}
