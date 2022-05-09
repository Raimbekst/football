package repository

import (
	"carWash/internal/domain"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"golang.org/x/net/context"
	"time"
)

type FavouriteRepos struct {
	db *sqlx.DB
}

func NewFavouriteRepos(db *sqlx.DB) *FavouriteRepos {
	return &FavouriteRepos{db: db}
}

func (f *FavouriteRepos) Create(ctx *fiber.Ctx, input FavouriteInput) (int, error) {
	var id int

	_, cancel := context.WithTimeout(ctx.Context(), time.Second*4)

	defer cancel()

	query := fmt.Sprintf("INSERT INTO %s(user_id,building_id) VALUES($1,$2) RETURNING id", favouriteTable)

	err := f.db.QueryRowx(query, input.UserId, input.BuildingId).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("repository.Create: %w", err)
	}

	return id, nil
}

func (f *FavouriteRepos) GetAll(ctx *fiber.Ctx, page domain.Pagination, userId int) (*domain.GetAllResponses, error) {
	var (
		setValues string
	)

	_, cancel := context.WithTimeout(ctx.Context(), 4*time.Second)

	defer cancel()

	setValues = fmt.Sprintf(" WHERE user_id = %d", userId)

	count, err := countPage(f.db, favouriteTable, setValues)

	if err != nil {
		return nil, fmt.Errorf("repository.GetAll: %w", err)
	}

	offset, pagesCount := calculatePagination(&page, count)

	inp := make([]*domain.Favourite, 0, page.Limit)

	query := fmt.Sprintf(
		`SELECT 
					f.id,
					bu.id "b.id",
					bu.building_name "b.building_name",
					bu.address "b.address",
					bu.instagram "b.instagram", 
					bu.description "b.description" 
				FROM 
					%s f
				INNER JOIN 
					%s bu
				ON 
					f.building_id = bu.id
				%s
				    ORDER BY
				f.id ASC LIMIT $1 OFFSET $2`, favouriteTable, buildingTable, setValues)

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

func (f *FavouriteRepos) GetById(ctx *fiber.Ctx, id, userId int) (*domain.Favourite, error) {
	var (
		inp domain.Favourite
	)

	_, cancel := context.WithTimeout(ctx.Context(), 4*time.Second)

	defer cancel()

	query := fmt.Sprintf(
		`SELECT 
					f.id "id",
					bu.id "b.id",
					bu.building_name "b.building_name",
					bu.address "b.address",
					bu.instagram "b.instagram",
					bu.description "b.description" 
				FROM 
					%s f
				INNER JOIN 
					%s bu
				ON 
					f.building_id = bu.id

				WHERE f.id = $1 AND f.user_id = $2`, favouriteTable, buildingTable)

	err := f.db.Get(&inp, query, id, userId)

	if err != nil {
		return nil, fmt.Errorf("repository.GetById: %w", domain.ErrNotFound)
	}

	return &inp, nil
}

func (f *FavouriteRepos) Delete(ctx *fiber.Ctx, id, userId int) error {

	_, cancel := context.WithTimeout(ctx.Context(), 4*time.Second)

	defer cancel()

	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1 AND user_id = $2", favouriteTable)

	result, err := f.db.Exec(query, id, userId)

	affected, err := result.RowsAffected()

	if err != nil {
		return fmt.Errorf("repository.Delete: %w", err)
	}

	if affected == 0 {
		return fmt.Errorf("repository.Delete: %w", domain.ErrNotFound)
	}

	return nil
}
