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

type PitchRepos struct {
	db *sqlx.DB
}

func NewPitchRepos(db *sqlx.DB) *PitchRepos {
	return &PitchRepos{db: db}
}

func (p *PitchRepos) Create(ctx *fiber.Ctx, pitch domain.Pitch) (int, error) {
	var id int

	_, cancel := context.WithTimeout(ctx.Context(), 4*time.Second)

	defer cancel()

	query := fmt.Sprintf(
		`INSERT 
					INTO 
					%s 
						(building_id,price,pitch_image) 
					SELECT 
						b.id, $1, $2 
					FROM %s b INNER JOIN %s u on b.manager_id = u.id WHERE b.id = $3 AND u.id = $4 RETURNING id`, pitchTable, buildingTable, userTable)

	err := p.db.QueryRowx(query, pitch.Price, pitch.Image, pitch.BuildingId, pitch.ManagerId).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("repository.Create: %w", err)
	}
	return id, nil
}

func (p *PitchRepos) GetAll(ctx *fiber.Ctx, page domain.Pagination, id int) (*domain.GetAllResponses, error) {

	var (
		setValues      string
		forCheckValues []string
		whereClause    string
		url            = ctx.BaseURL()
	)

	_, cancel := context.WithTimeout(ctx.Context(), 4*time.Second)

	defer cancel()

	if id != 0 {
		forCheckValues = append(forCheckValues, fmt.Sprintf("building_id = %d", id))
	}

	whereClause = strings.Join(forCheckValues, " AND ")

	if whereClause != "" {
		setValues = "WHERE " + whereClause
	}

	count, err := countPage(p.db, pitchTable, setValues)
	if err != nil {
		return nil, fmt.Errorf("repository.GetAll: %w", err)
	}

	offset, pagesCount := calculatePagination(&page, count)

	inp := make([]*domain.Pitch, 0, page.Limit)

	query := fmt.Sprintf("SELECT * FROM %s ORDER BY id ASC LIMIT $1 OFFSET $2", pitchTable)

	err = p.db.Select(&inp, query, page.Limit, offset)

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

func (p *PitchRepos) GetById(ctx *fiber.Ctx, id int) (*domain.Pitch, error) {
	var (
		place domain.Pitch
		url   = ctx.BaseURL()
	)

	_, cancel := context.WithTimeout(ctx.Context(), 4*time.Second)

	defer cancel()

	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", pitchTable)

	err := p.db.Get(&place, query, id)

	if err != nil {
		return nil, fmt.Errorf("repository.GetById: %w", domain.ErrNotFound)
	}
	place.Image = url + "/" + "media/" + place.Image

	return &place, nil
}

func (p *PitchRepos) Update(ctx *fiber.Ctx, id int, inp domain.Pitch) ([]string, error) {

	_, cancel := context.WithTimeout(ctx.Context(), 4*time.Second)

	defer cancel()

	setValues := make([]string, 0, reflect.TypeOf(domain.Pitch{}).NumField())

	var images []string

	if inp.BuildingId != 0 {
		setValues = append(setValues, fmt.Sprintf("building_id=:building_id"))
	}

	if inp.Image != "" {
		setValues = append(setValues, fmt.Sprintf("pitch_image=:pitch_image"))
		images = append(images, "pitch_image")
	}
	setImages := strings.Join(images, ", ")

	var input domain.Pitch

	querySelectImages := fmt.Sprintf("SELECT %s FROM %s WHERE id = $1", setImages, pitchTable)

	err := p.db.Get(&input, querySelectImages, id)

	images = nil

	images = append(images, input.Image)

	setQuery := strings.Join(setValues, ", ")

	if setQuery == "" {
		return nil, fmt.Errorf("repository.Update: %w", errors.New("empty body"))
	}

	query := fmt.Sprintf(`UPDATE %s p SET %s FROM %s b, %s u WHERE p.building_id = b.id AND b.manager_id = u.id AND p.id = %d`, pitchTable, setQuery, buildingTable, userTable, id)

	result, err := p.db.NamedExec(query, inp)

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

func (p *PitchRepos) Delete(ctx *fiber.Ctx, id int) ([]string, error) {

	var (
		image  string
		images []string
	)

	_, cancel := context.WithTimeout(ctx.Context(), 500*time.Millisecond)

	defer cancel()

	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1 RETURNING pitch_image", pitchTable)

	err := p.db.QueryRowx(query, id).Scan(&image)

	if err != nil {
		return nil, fmt.Errorf("repository.Delete: %w", domain.ErrNotFound)
	}

	images = append(images, image)

	return images, nil
}
