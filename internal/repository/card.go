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

type CardRepos struct {
	db *sqlx.DB
}

func NewCardRepos(db *sqlx.DB) *CardRepos {
	return &CardRepos{db: db}
}

func (f *CardRepos) Create(ctx *fiber.Ctx, card domain.Card) (int, error) {
	var id int

	_, cancel := context.WithTimeout(ctx.Context(), 4*time.Second)

	defer cancel()

	query := fmt.Sprintf(
		`INSERT 
					INTO 
					%s 
						(user_id,cvv,full_name,full_number)
					VALUES
						($1,$2,$3,$4)
			 RETURNING id`, cardTable)

	err := f.db.QueryRowx(query, card.UserId, card.Cvv, card.FullName, card.FullNumber).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("repository.Create: %w", err)
	}
	return id, nil
}

func (f *CardRepos) GetAll(ctx *fiber.Ctx, page domain.Pagination, userId int) (*domain.GetAllResponses, error) {

	_, cancel := context.WithTimeout(ctx.Context(), 4*time.Second)

	defer cancel()

	count, err := countPage(f.db, cardTable, "")
	if err != nil {
		return nil, fmt.Errorf("repository.GetAll: %w", err)
	}

	offset, pagesCount := calculatePagination(&page, count)

	inp := make([]*domain.Card, 0, page.Limit)

	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1 ORDER BY id ASC LIMIT $2 OFFSET $3", cardTable)

	err = f.db.Select(&inp, query, userId, page.Limit, offset)

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

func (f *CardRepos) GetById(ctx *fiber.Ctx, id, userId int) (*domain.Card, error) {
	var (
		place domain.Card
	)

	_, cancel := context.WithTimeout(ctx.Context(), 4*time.Second)

	defer cancel()

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1 AND user_id = $2", cardTable)

	err := f.db.Get(&place, query, id, userId)

	if err != nil {
		return nil, fmt.Errorf("repository.GetById: %w", domain.ErrNotFound)
	}

	return &place, nil
}

func (f *CardRepos) Update(ctx *fiber.Ctx, id int, inp domain.Card) error {
	_, cancel := context.WithTimeout(ctx.Context(), 4*time.Second)

	defer cancel()

	setValues := make([]string, 0, reflect.TypeOf(domain.Card{}).NumField())

	if inp.Cvv != 0 {
		setValues = append(setValues, fmt.Sprintf("cvv=:cvv"))
	}

	if inp.FullName != "" {
		setValues = append(setValues, fmt.Sprintf("full_name=:full_name"))
	}

	if inp.FullNumber != "" {
		setValues = append(setValues, fmt.Sprintf("full_number=:full_number"))
	}

	setQuery := strings.Join(setValues, ", ")

	if setQuery == "" {
		return fmt.Errorf("repository.Update: %w", errors.New("empty body"))
	}

	query := fmt.Sprintf(`UPDATE %s SET %s WHERE id = %d AND user_id = %d`, cardTable, setQuery, id, inp.UserId)

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

func (f *CardRepos) Delete(ctx *fiber.Ctx, id, userId int) error {

	_, cancel := context.WithTimeout(ctx.Context(), 500*time.Millisecond)

	defer cancel()

	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1 AND user_id = $2 ", cardTable)

	err := f.db.QueryRowx(query, id, userId).Err()

	if err != nil {
		return fmt.Errorf("repository.Delete: %w", domain.ErrNotFound)
	}
	return nil

}
