package repository

import (
	"carWash/internal/domain"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"time"
)

type FeedbackRepos struct {
	db *sqlx.DB
}

func (f *FeedbackRepos) CreateNoty(ctx *fiber.Ctx, noty domain.Notification) (int, error) {
	var id int

	_, cancel := context.WithTimeout(ctx.Context(), 4*time.Second)

	defer cancel()

	query := fmt.Sprintf("INSERT INTO %s(title,content) VALUES($1,$2) RETURNING id", notificationTable)

	err := f.db.QueryRowx(query, noty.Title, noty.Content).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("repository.Create: %w", err)
	}
	return id, nil
}

func (f *FeedbackRepos) GetAllNoty(ctx *fiber.Ctx, page domain.Pagination) (*domain.GetAllResponses, error) {

	var setValues string

	_, cancel := context.WithTimeout(ctx.Context(), 500*time.Millisecond)

	defer cancel()

	count, err := countPage(f.db, notificationTable, setValues)

	if err != nil {
		return nil, fmt.Errorf("repository.GetAll : %w", err)
	}

	offset, pagesCount := calculatePagination(&page, count)

	inp := make([]*domain.Notification, 0, page.Limit)

	query := fmt.Sprintf(
		`SELECT
					id,title,content
				FROM 
					%s 
					ORDER BY 
				id ASC LIMIT $1 OFFSET $2`, notificationTable)

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

func NewFeedbackRepos(db *sqlx.DB) *FeedbackRepos {
	return &FeedbackRepos{db: db}
}

func (f *FeedbackRepos) Create(ctx *fiber.Ctx, feedback domain.Feedback, userId int) (int, error) {
	var id int

	_, cancel := context.WithTimeout(ctx.Context(), 4*time.Second)

	defer cancel()

	query := fmt.Sprintf("INSERT INTO %s(user_id,text) VALUES($1,$2) RETURNING id", feedbackTable)

	err := f.db.QueryRowx(query, userId, feedback.Text).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("repository.Create: %w", err)
	}
	return id, nil
}

func (f *FeedbackRepos) GetAll(ctx *fiber.Ctx, page domain.Pagination) (*domain.GetAllResponses, error) {
	var setValues string

	_, cancel := context.WithTimeout(ctx.Context(), 500*time.Millisecond)

	defer cancel()

	count, err := countPage(f.db, feedbackTable, setValues)

	if err != nil {
		return nil, fmt.Errorf("repository.GetAll : %w", err)
	}

	offset, pagesCount := calculatePagination(&page, count)

	inp := make([]*domain.Feedback, 0, page.Limit)

	query := fmt.Sprintf(
		`SELECT
					f.id,
					f.user_id,
					f.text,
					u.phone_number,
					u.user_name	
				FROM 
					%s f
				INNER JOIN 
					%s u
				ON 
				f.user_id = u.id
					ORDER BY 
				f.id ASC LIMIT $1 OFFSET $2`, feedbackTable, userTable)

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
