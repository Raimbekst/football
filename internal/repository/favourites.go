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
		count int
		url   = ctx.BaseURL()
	)

	_, cancel := context.WithTimeout(ctx.Context(), 4*time.Second)

	defer cancel()

	inp := make([]*domain.Building, 0)

	query := fmt.Sprintf(
		`SELECT
					b.*,
					u.phone_number,
					COALESCE(min(price),null,0) as min_price,
					f.id  "f.id"
				FROM 
					%s b
				LEFT OUTER JOIN 
					%s p
				ON 
					b.id = p.building_id
				LEFT OUTER JOIN 
					%s f
				ON 
					b.id = f.building_id
				LEFT OUTER JOIN 
					%s u 
				ON 
					b.manager_id = u.id
				WHERE 
					f.user_id = $1
				group by
    				b.id,b.building_image, building_name,
    				address, instagram,
					manager_id, description,
    				work_time, start_time,
    				end_time, longtitude, latitude,f.id, f.user_id, u.phone_number`, buildingTable, pitchTable, favouriteTable, userTable)

	err := f.db.Select(&inp, query, userId)

	if err != nil {
		return nil, fmt.Errorf("repository.GetAll: %w", err)
	}

	for _, val := range inp {
		val.BuildingImage = url + "/" + "media/" + val.BuildingImage
	}

	pages := domain.PaginationPage{
		Page:  page.Page,
		Count: count,
	}

	ans := domain.GetAllResponses{
		Data:     inp,
		PageInfo: pages,
	}
	return &ans, nil
}

func (f *FavouriteRepos) GetById(ctx *fiber.Ctx, id, userId int) (*domain.Building, error) {
	var (
		inp domain.Building
	)

	_, cancel := context.WithTimeout(ctx.Context(), 4*time.Second)

	defer cancel()

	query := fmt.Sprintf(
		`SELECT
					b.*,
					u.phone_number,
					COALESCE(min(price),null,0) as min_price,
					f.id  "f.id"
				FROM 
					%s f
				INNER 	JOIN 
					%s b
				ON 
					b.id = f.building_id
				LEFT OUTER JOIN 
					%s p
				ON 
					b.id = p.building_id
				LEFT OUTER JOIN 
					%s u 
				ON 
					b.manager_id = u.id
				WHERE 
					f.user_id = $1 AND f.building_id = $2
				group by
    				b.id,b.building_image, building_name,
    				address, instagram,
					manager_id, description,
    				work_time, start_time,
    				end_time, longtitude, latitude,f.id, f.user_id, u.phone_number`, favouriteTable, buildingTable, pitchTable, userTable)

	err := f.db.Get(&inp, query, userId, id)

	if err != nil {

		return nil, fmt.Errorf("repository.GetById: %w", domain.ErrNotFound)
	}

	return &inp, nil
}

func (f *FavouriteRepos) Delete(ctx *fiber.Ctx, id, userId int) error {

	_, cancel := context.WithTimeout(ctx.Context(), 4*time.Second)

	defer cancel()

	query := fmt.Sprintf("DELETE FROM %s WHERE building_id = $1 AND user_id = $2", favouriteTable)

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
