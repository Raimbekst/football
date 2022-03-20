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

type BuildingRepos struct {
	db *sqlx.DB
}

func NewBuildingRepos(db *sqlx.DB) *BuildingRepos {
	return &BuildingRepos{db: db}
}

func (b *BuildingRepos) Create(ctx *fiber.Ctx, building domain.Building) (int, error) {
	var id int

	_, cancel := context.WithTimeout(ctx.Context(), 4*time.Second)

	defer cancel()

	query := fmt.Sprintf("INSERT INTO %s(building_name, address, instagram, manager_id) VALUES($1,$2,$3,$4) RETURNING id", buildingTable)

	err := b.db.QueryRowx(query, building.Name, building.Address, building.Instagram, building.ManagerId).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("repository.Create: %w", err)
	}
	return id, nil
}

func (b *BuildingRepos) GetAll(ctx *fiber.Ctx, page domain.Pagination, info domain.UserInfo) (*domain.GetAllResponses, error) {

	var (
		setValues string
	)

	_, cancel := context.WithTimeout(ctx.Context(), 4*time.Second)

	defer cancel()

	switch info.Type {

	case "manager":
		setValues = fmt.Sprintf("WHERE manager_id = %d", info.Id)
	}

	count, err := countPage(b.db, buildingTable, setValues)
	if err != nil {
		return nil, fmt.Errorf("repository.GetAll: %w", err)
	}

	offset, pagesCount := calculatePagination(&page, count)

	inp := make([]*domain.Building, 0, page.Limit)

	query := fmt.Sprintf("SELECT * FROM %s %s ORDER BY id ASC LIMIT $1 OFFSET $2", buildingTable, setValues)

	err = b.db.Select(&inp, query, page.Limit, offset)

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

func (b *BuildingRepos) GetById(ctx *fiber.Ctx, id int) (*domain.Building, error) {

	var (
		inp domain.Building
	)

	_, cancel := context.WithTimeout(ctx.Context(), 4*time.Second)

	defer cancel()

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", buildingTable)

	err := b.db.Get(&inp, query, id)

	if err != nil {
		return nil, fmt.Errorf("repository.GetById: %w", domain.ErrNotFound)
	}

	return &inp, nil
}

func (b *BuildingRepos) Update(ctx *fiber.Ctx, id int, inp domain.Building) error {

	setValues := make([]string, 0, reflect.TypeOf(domain.Building{}).NumField())

	if inp.Name != "" {
		setValues = append(setValues, fmt.Sprintf("building_name=:building_name"))
	}

	if inp.Address != "" {
		setValues = append(setValues, fmt.Sprintf("address=:address"))
	}

	if inp.Instagram != "" {
		setValues = append(setValues, fmt.Sprintf("instagram=:instagram"))
	}

	_, cancel := context.WithTimeout(ctx.Context(), 500*time.Millisecond)

	defer cancel()

	setQuery := strings.Join(setValues, ", ")

	if setQuery == "" {
		return fmt.Errorf("repository.Update: %w", errors.New("empty body"))
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = %d AND manager_id = %d", buildingTable, setQuery, id, inp.ManagerId)

	result, err := b.db.NamedExec(query, inp)

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
func (b *BuildingRepos) Delete(ctx *fiber.Ctx, id int) error {
	_, cancel := context.WithTimeout(ctx.Context(), 4*time.Second)

	defer cancel()

	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", buildingTable)

	result, err := b.db.Exec(query, id)

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository.Delete: %w", err)
	}

	if affected == 0 {
		return fmt.Errorf("repository.Delete: %w", domain.ErrNotFound)
	}

	return nil
}

func countPage(db *sqlx.DB, table, setValues string) (int, error) {

	var count int
	queryCount := fmt.Sprintf("SELECT COUNT(*) FROM %s %s", table, setValues)

	row := db.QueryRow(queryCount)
	err := row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("repository.countPage : %w", err)
	}
	return count, nil
}
