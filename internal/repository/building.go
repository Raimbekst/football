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

	tx := b.db.MustBegin()

	_, cancel := context.WithTimeout(ctx.Context(), 4*time.Second)

	defer cancel()

	query := fmt.Sprintf("INSERT INTO %s(building_name, address, instagram, manager_id,description,building_image,work_time,start_time,end_time,longtitude,latitude) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING id", buildingTable)

	err := tx.QueryRowx(query, building.Name, building.Address, building.Instagram, building.ManagerId, building.Description, building.BuildingImage, building.WorkTime, secondToTime(building.StartTime), secondToTime(building.EndTime), building.Longtitude, building.Latitude).Scan(&id)

	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			return 0, fmt.Errorf("repository.Create: %w", txErr)
		}
		return 0, fmt.Errorf("repository.Create: %w", err)
	}

	for i := building.StartTime; i <= building.EndTime; i = i + 1800 {
		query := fmt.Sprintf("INSERT INTO %s(work_time, building_id) VALUES($1,$2)", timeTable)
		_, err := tx.Exec(query, secondToTime(i), id)
		if err != nil {
			txErr := tx.Rollback()
			if txErr != nil {
				return 0, fmt.Errorf("repository.Create: %w", txErr)
			}
			return 0, fmt.Errorf("repository.Create: %w", err)
		}
	}

	txErr := tx.Commit()
	if txErr != nil {
		return 0, fmt.Errorf("repository.Create: %w", txErr)
	}

	return id, nil
}

func (b *BuildingRepos) GetAll(ctx *fiber.Ctx, page domain.Pagination, info domain.UserInfo, building domain.FilterForBuilding) (*domain.GetAllResponses, error) {

	var (
		countValues      string
		setValues        string
		whereClause      string = " WHERE "
		havingClause     string = " Having "
		whereValuesList  []string
		havingValuesList []string
		count            int
		userId           int
		userIdCase       string
		url              = ctx.BaseURL()
	)

	_, cancel := context.WithTimeout(ctx.Context(), 4*time.Second)

	defer cancel()

	switch info.Type {

	case "manager":
		whereValuesList = append(whereValuesList, fmt.Sprintf("manager_id = %d", info.Id))
	case "user":
		userIdCase = fmt.Sprintf(
			`,CASE 
						WHEN f.user_id = %d THEN true 
						ELSE false 
					END is_favourite`, info.Id)
		userId = info.Id
	}

	if building.PitchType != 0 {
		whereValuesList = append(whereValuesList, fmt.Sprintf("pitch_type = %d", building.PitchType))
	}

	if building.PitchExtra != 0 {
		whereValuesList = append(whereValuesList, fmt.Sprintf("pitch_extra = %d", building.PitchExtra))
	}

	if building.StartCost != nil {
		havingValuesList = append(havingValuesList, fmt.Sprintf("min(price) >= %d", *building.StartCost))
	}

	if building.EndCost != nil {
		havingValuesList = append(havingValuesList, fmt.Sprintf("min(price) <= %d", *building.EndCost))
	}

	whereValuesJoin := strings.Join(whereValuesList, " AND ")

	if whereValuesList != nil {
		countValues = countValues + whereClause + whereValuesJoin
		setValues = setValues + whereClause + whereValuesJoin
	}

	setValues = setValues + fmt.Sprintf(" GROUP BY b.id,b.building_image, building_name,building_image, address, instagram, manager_id, description, work_time, start_time, end_time,longtitude,latitude,f.user_id,u.phone_number")

	havingValuesJoin := strings.Join(havingValuesList, " AND ")

	if havingValuesList != nil {
		havingValuesJoin = havingClause + havingValuesJoin
		setValues = setValues + havingValuesJoin
	}

	queryCount := fmt.Sprintf(
		`WITH new_building AS 
					(SELECT 
							count(distinct b.id) 
						FROM
							%s b
						LEFT OUTER  JOIN 
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
						%s
						GROUP BY b.id
					%s)
			 SELECT COUNT(*) FROM new_building ;`, buildingTable, pitchTable, favouriteTable, userTable, countValues, havingValuesJoin)

	_ = b.db.QueryRowx(queryCount).Scan(&count)

	offset, pagesCount := calculatePagination(&page, count)

	inp := make([]*domain.Building, 0)

	query := fmt.Sprintf(
		`SELECT 
					b.*,
					u.phone_number,
					COALESCE(min(price),null,0) as min_price
					%s	
				FROM 
					%s b
				LEFT OUTER JOIN 
					%s p
				ON 
					b.id = p.building_id
				LEFT OUTER JOIN 
					%s f
				ON 
					(b.id = f.building_id AND f.user_id = %d)
				LEFT OUTER JOIN 
					%s u 
				ON 
					b.manager_id = u.id
				%s
				ORDER BY
					b.id ASC LIMIT $1 OFFSET $2`, userIdCase, buildingTable, pitchTable, favouriteTable, userId, userTable, setValues)

	err := b.db.Select(&inp, query, page.Limit, offset)
	fmt.Println(query)
	if err != nil {
		return nil, fmt.Errorf("repository.GetAll: %w", err)
	}

	for _, val := range inp {
		val.BuildingImage = url + "/" + "media/" + val.BuildingImage

		query := fmt.Sprintf(
			`SELECT
					coalesce(AVG(grade),null,0) 
				FROM 
					%s WHERE building_id = $1`, gradeTable)

		row := b.db.QueryRow(query, val.Id)
		err = row.Scan(&val.Grade)

		if err != nil {
			return nil, fmt.Errorf("repository.GetAll: %w", err)
		}

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

func (b *BuildingRepos) GetById(ctx *fiber.Ctx, info domain.UserInfo, id int) (*domain.Building, error) {

	var (
		inp        domain.Building
		url        = ctx.BaseURL()
		userIdCase string
	)

	_, cancel := context.WithTimeout(ctx.Context(), 4*time.Second)

	defer cancel()

	if info.Type == "user" {
		userIdCase = fmt.Sprintf(
			`,CASE 
						WHEN f.user_id = %d THEN true 
						ELSE false 
					END is_favourite`, info.Id)
	}

	query := fmt.Sprintf(`
				SELECT 
					b.*,
					u.phone_number,	
					COALESCE(min(price),null,0) as min_price
					%s
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
					b.id = $1
				GROUP BY 
					b.id, building_name,building_image, address, instagram, manager_id, description, work_time, start_time, end_time, longtitude, latitude,f.user_id,u.phone_number;`, userIdCase, buildingTable, pitchTable, favouriteTable, userTable)

	err := b.db.Get(&inp, query, id)

	inp.BuildingImage = url + "/" + "media/" + inp.BuildingImage

	if err != nil {
		return nil, fmt.Errorf("repository.GetById: %w", domain.ErrNotFound)
	}

	return &inp, nil
}

func (b *BuildingRepos) Update(ctx *fiber.Ctx, id int, inp domain.Building) ([]string, error) {

	setValues := make([]string, 0, reflect.TypeOf(domain.Building{}).NumField())

	var images []string

	if inp.Name != "" {
		setValues = append(setValues, fmt.Sprintf("building_name=:building_name"))
	}

	if inp.Address != "" {
		setValues = append(setValues, fmt.Sprintf("address=:address"))
	}

	if inp.Instagram != "" {
		setValues = append(setValues, fmt.Sprintf("instagram=:instagram"))
	}
	if inp.BuildingImage != "" {
		setValues = append(setValues, fmt.Sprintf("building_image=:building_image"))
		images = append(images, "building_image")
	}

	_, cancel := context.WithTimeout(ctx.Context(), 4*time.Second)

	defer cancel()

	setImages := strings.Join(images, ", ")

	var input domain.Building

	querySelectImages := fmt.Sprintf("SELECT %s FROM %s WHERE id = $1", setImages, buildingTable)

	err := b.db.Get(&input, querySelectImages, id)

	images = nil

	images = append(images, input.BuildingImage)

	setQuery := strings.Join(setValues, ", ")

	if setQuery == "" {
		return nil, fmt.Errorf("repository.Update: %w", errors.New("empty body"))
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = %d AND manager_id = %d", buildingTable, setQuery, id, inp.ManagerId)

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

func (b *BuildingRepos) Delete(ctx *fiber.Ctx, id int) ([]string, error) {
	var (
		image  string
		images []string
	)

	_, cancel := context.WithTimeout(ctx.Context(), 4*time.Second)

	defer cancel()

	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1 RETURNING building_image", buildingTable)

	err := b.db.QueryRowx(query, id).Scan(&image)

	if err != nil {
		return nil, fmt.Errorf("repository.Delete: %w", err)
	}
	images = append(images, image)
	return images, nil
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

func secondToTime(a int) string {
	min := a / 60

	a = a % 60

	hh := min / 60

	min = min % 60
	return fmt.Sprintf("%d:%d", hh, min)
}
