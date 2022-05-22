package repository

import (
	"carWash/internal/domain"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"time"
)

type CommentRepos struct {
	db *sqlx.DB
}

func NewCommentRepos(db *sqlx.DB) *CommentRepos {
	return &CommentRepos{db: db}
}

func (c *CommentRepos) Create(ctx *fiber.Ctx, comment domain.Comment) (int, error) {

	var id int
	var inp domain.Comment

	_, cancel := context.WithTimeout(ctx.Context(), 500*time.Millisecond)

	defer cancel()

	queryCheck := fmt.Sprintf(
		`SELECT * FROM %s WHERE user_id = $1 AND building_id = $2`, commentTable)

	err := c.db.Get(&inp, queryCheck, comment.UserId, comment.BuildingId)

	if err == nil {
		return 0, fmt.Errorf("repository.Create: %w", domain.ErrUserCommented)
	}

	query := fmt.Sprintf(
		`INSERT INTO %s(comment,user_id,building_id) VALUES($1,$2,$3) RETURNING id`, commentTable)

	err = c.db.QueryRowx(query, comment.CommentText, comment.UserId, comment.BuildingId).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("repository.Create: %w", err)
	}
	return id, nil
}

func (c *CommentRepos) GetAll(ctx *fiber.Ctx, page domain.Pagination, buildingId int) (*domain.GetAllResponses, error) {

	var setValues string

	_, cancel := context.WithTimeout(ctx.Context(), 500*time.Millisecond)
	defer cancel()

	if buildingId != 0 {
		setValues = fmt.Sprintf("WHERE building_id = %d", buildingId)
	}

	count, err := countPage(c.db, commentTable, setValues)

	if err != nil {
		return nil, fmt.Errorf("repository.GetAll : %w", err)
	}

	offset, pagesCount := calculatePagination(&page, count)

	inp := make([]*domain.Comment, 0, page.Limit)

	query := fmt.Sprintf(
		`SELECT
					com.id,u.user_name,building_id,comment
				FROM 
					%s com  
				INNER JOIN 
					%s u
				ON 
					com.user_id = u.id
				%s 
					ORDER BY 
				com.id ASC LIMIT $1 OFFSET $2`, commentTable, userTable, setValues)

	err = c.db.Select(&inp, query, page.Limit, offset)

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

func (c *CommentRepos) CreateGrade(ctx *fiber.Ctx, grade domain.Grade) (int, error) {
	var id int
	var inp domain.Comment

	_, cancel := context.WithTimeout(ctx.Context(), 500*time.Millisecond)

	defer cancel()

	queryCheck := fmt.Sprintf(
		`SELECT * FROM %s WHERE user_id = $1 AND building_id = $2`, gradeTable)

	err := c.db.Get(&inp, queryCheck, grade.UserId, grade.BuildingId)

	if err == nil {
		return 0, fmt.Errorf("repository.Create: %w", domain.ErrUserCommented)
	}

	query := fmt.Sprintf(
		`INSERT INTO %s(grade,user_id,building_id) VALUES($1,$2,$3) RETURNING id`, gradeTable)

	err = c.db.QueryRowx(query, grade.Grade, grade.UserId, grade.BuildingId).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("repository.Create: %w", err)
	}
	return id, nil

}

func (c *CommentRepos) GetAllGrades(ctx *fiber.Ctx, page domain.Pagination, buildingId int) (*domain.GetAllResponses, error) {
	var setValues string

	_, cancel := context.WithTimeout(ctx.Context(), 500*time.Millisecond)
	defer cancel()

	if buildingId != 0 {
		setValues = fmt.Sprintf("WHERE building_id = %d", buildingId)
	}
	var grade float64

	query := fmt.Sprintf(
		`SELECT
					coalesce(AVG(grade),null,0) 
				FROM 
					%s %s`, gradeTable, setValues)

	row := c.db.QueryRow(query)
	err := row.Scan(&grade)

	if err != nil {
		return nil, fmt.Errorf("repository.GetAll: %w", err)
	}

	pages := domain.PaginationPage{
		Page: page.Page,
	}
	ans := domain.GetAllResponses{
		Data:     grade,
		PageInfo: pages,
	}
	return &ans, nil
}
