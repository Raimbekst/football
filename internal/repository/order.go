package repository

import (
	"carWash/internal/domain"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
)

type OrderRepos struct {
	db *sqlx.DB
}

func NewOrderRepos(db *sqlx.DB) *OrderRepos {
	return &OrderRepos{db: db}
}

func (o *OrderRepos) Create(ctx *fiber.Ctx, order domain.Order) (int, error) {
	var id int

	_, cancel := context.WithTimeout(ctx.Context(), 4*time.Second)

	defer cancel()

	query := fmt.Sprintf(
		`INSERT INTO
					%s
				(pitch_id,user_id,order_date,start_time,end_time,status) 
					VALUES
				($1,$2,to_timestamp($3) at time zone 'GMT',$4,$5,$6) RETURNING id`, orderTable)

	err := o.db.QueryRowx(query, order.PitchId, order.UserId, order.OrderDate, order.StartTime, order.EndTime, order.Status).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("repository.Create: %w", err)
	}
	return id, nil
}

func (o *OrderRepos) GetAll(ctx *fiber.Ctx, page domain.Pagination, info domain.UserInfo, date float64) (*domain.GetAllResponses, error) {
	var (
		setValues      string
		forCheckValues []string
		whereClause    string
		url            = ctx.BaseURL()
		count          int
	)

	_, cancel := context.WithTimeout(ctx.Context(), 4*time.Second)

	defer cancel()

	switch info.Type {
	case "manager":
		forCheckValues = append(forCheckValues, fmt.Sprintf("b.manager_id = %d", info.Id))

	case "user":
		forCheckValues = append(forCheckValues, fmt.Sprintf("o.user_id = %d", info.Id))
	}

	if date != 0 {
		forCheckValues = append(forCheckValues, fmt.Sprintf("o.order_date = to_timestamp(%f) at time zone 'GMT'", date))
	}

	whereClause = strings.Join(forCheckValues, " AND ")

	if whereClause != "" {
		setValues = "WHERE " + whereClause
	}

	queryCount := fmt.Sprintf(
		`select 
					count(*)
				from
					%s o 
				LEFT OUTER JOIN
					%s p 
				on 
					p.id = o.pitch_id 
				LEFT OUTER JOIN
					%s b 
				on 
					b.id = p.building_id
				LEFT OUTER JOIN
					%s u on o.user_id = u.id
				%s`, orderTable, pitchTable, buildingTable, userTable, setValues)

	err := o.db.QueryRowx(queryCount).Scan(&count)

	if err != nil {
		return nil, fmt.Errorf("repository.GetAll: %w", err)
	}

	offset, pagesCount := calculatePagination(&page, count)

	inp := make([]*domain.Order, 0, page.Limit)

	query := fmt.Sprintf(
		`select 
					o.id,
					extract(epoch from o.order_date::timestamp at time zone 'GMT') "order_date",
					o.start_time,
					o.end_time,
					o.status,
					p.price,
					p.pitch_type,
					p.pitch_image,
					b.building_name,
					b.address,
					u.user_name,
					u.phone_number
				from 
					%s o 
				LEFT OUTER JOIN
					%s p 
				on 
					p.id = o.pitch_id
				LEFT OUTER JOIN
					%s b
				on 
					b.id = p.building_id
				LEFT OUTER JOIN
					%s u
				on 
					o.user_id = u.id 
				%s ORDER BY o.id ASC LIMIT $1 OFFSET $2`, orderTable, pitchTable, buildingTable, userTable, setValues)

	err = o.db.Select(&inp, query, page.Limit, offset)

	if err != nil {
		return nil, fmt.Errorf("repository.GetAll: %w", err)
	}

	pages := domain.PaginationPage{
		Page:  page.Page,
		Pages: pagesCount,
		Count: count,
	}

	for _, value := range inp {
		value.PitchImage = url + "/" + "media/" + value.PitchImage
	}

	ans := domain.GetAllResponses{
		Data:     inp,
		PageInfo: pages,
	}
	return &ans, nil
}
