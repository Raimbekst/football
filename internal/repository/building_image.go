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

type BuildingImageRepos struct {
	db *sqlx.DB
}

func NewBuildingImageRepos(db *sqlx.DB) *BuildingImageRepos {
	return &BuildingImageRepos{db: db}
}

func (b *BuildingImageRepos) Create(ctx *fiber.Ctx, building domain.BuildingImage) (int, error) {
	var id int

	_, cancel := context.WithTimeout(ctx.Context(), 4*time.Second)

	defer cancel()

	query := fmt.Sprintf(`INSERT 	
									INTO 
								%s
									(building_id,building_image) 
								SELECT 
									b.id, $1 
								FROM 
									%s b INNER JOIN %s u 
								on b.manager_id = u.id 
									WHERE 
								b.id = $2 AND u.id = $3
								RETURNING id`, buildingImageTable, buildingTable, userTable)

	err := b.db.QueryRowx(query, building.Image, building.BuildingId, building.ManagerId).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("repository.Create: %w", err)
	}
	return id, nil

}

func (b *BuildingImageRepos) GetAll(ctx *fiber.Ctx, page domain.Pagination, id int) (*domain.GetAllResponses, error) {
	var (
		setValues string
		url       = ctx.BaseURL()
	)

	_, cancel := context.WithTimeout(ctx.Context(), 4*time.Second)

	defer cancel()

	if id != 0 {
		setValues = fmt.Sprintf("WHERE building_id = %d", id)
	}

	count, err := countPage(b.db, buildingImageTable, setValues)

	if err != nil {
		return nil, fmt.Errorf("repository.GetAll: %w", err)
	}

	offset, pagesCount := calculatePagination(&page, count)

	inp := make([]*domain.BuildingImage, 0, page.Limit)

	query := fmt.Sprintf("SELECT * FROM %s %s ORDER BY id ASC LIMIT $1 OFFSET $2", buildingImageTable, setValues)

	err = b.db.Select(&inp, query, page.Limit, offset)

	if err != nil {
		return nil, fmt.Errorf("repository.GetAll: %w", err)
	}

	pages := domain.PaginationPage{
		Page:  page.Page,
		Pages: pagesCount,
		Count: count,
	}

	for _, value := range inp {
		value.Image = url + "/" + "media/" + value.Image
	}

	ans := domain.GetAllResponses{
		Data:     inp,
		PageInfo: pages,
	}
	return &ans, nil
}

func (b BuildingImageRepos) GetById(ctx *fiber.Ctx, id int) (*domain.BuildingImage, error) {
	var (
		inp domain.BuildingImage
		url = ctx.BaseURL()
	)

	_, cancel := context.WithTimeout(ctx.Context(), 4*time.Second)

	defer cancel()

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", buildingImageTable)

	err := b.db.Get(&inp, query, id)

	if err != nil {
		return nil, fmt.Errorf("repository.GetById: %w", domain.ErrNotFound)
	}

	inp.Image = url + "/" + "media/" + inp.Image

	return &inp, nil
}

func (b BuildingImageRepos) Update(ctx *fiber.Ctx, id int, inp domain.BuildingImage) ([]string, error) {
	setValues := make([]string, 0, reflect.TypeOf(domain.BuildingImage{}).NumField())

	_, cancel := context.WithTimeout(ctx.Context(), 4*time.Second)

	defer cancel()

	var images []string

	if inp.BuildingId != 0 {
		setValues = append(setValues, fmt.Sprintf("building_id=:building_id"))
	}

	if inp.Image != "" {
		setValues = append(setValues, fmt.Sprintf("building_image=:building_image"))
		images = append(images, "building_image")
	}
	setImages := strings.Join(images, ", ")

	var input domain.BuildingImage

	querySelectImages := fmt.Sprintf("SELECT %s FROM %s WHERE id = $1", setImages, buildingImageTable)

	err := b.db.Get(&input, querySelectImages, id)

	images = nil

	images = append(images, input.Image)

	setQuery := strings.Join(setValues, ", ")

	if setQuery == "" {
		return nil, fmt.Errorf("repository.Update: %w", errors.New("empty body"))
	}

	query := fmt.Sprintf("UPDATE %s bi SET %s FROM %s b, %s u WHERE bi.building_id = b.id AND b.manager_id = u.id AND bi.id = %d", buildingImageTable, setQuery, buildingTable, userTable, id)

	result, err := b.db.NamedExec(query, inp)

	if err != nil {
		return nil, fmt.Errorf("repository.Update: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("repository.Update: %w", err)
	}

	if affected == 0 {
		return nil, fmt.Errorf("repository.Update: %w", domain.ErrNotFound)
	}

	return images, nil
}

func (b BuildingImageRepos) Delete(ctx *fiber.Ctx, id int) ([]string, error) {
	var (
		image  string
		images []string
	)

	_, cancel := context.WithTimeout(ctx.Context(), 500*time.Millisecond)

	defer cancel()

	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1 RETURNING building_image", buildingImageTable)

	err := b.db.QueryRowx(query, id).Scan(&image)

	if err != nil {
		return nil, fmt.Errorf("repository.Delete: %w", domain.ErrNotFound)
	}

	images = append(images, image)

	return images, nil
}
