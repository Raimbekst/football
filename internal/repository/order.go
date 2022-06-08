package repository

import (
	"carWash/internal/domain"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"strconv"
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

	tx := o.db.MustBegin()

	t := order.Times[(len(order.Times) - 1)]

	pt := strings.Split(t, ":")

	hh, err := strconv.Atoi(pt[0])

	if err != nil {
		return 0, fmt.Errorf("repository.Create:err1 %w", err)
	}

	min, err := strconv.Atoi(pt[1])

	if err != nil {
		return 0, fmt.Errorf("repository.Create:err1 %w", err)
	}

	total := hh*3600 + min*60

	query := fmt.Sprintf(
		`INSERT INTO
							%s
						(first_name, phone_number, extra_info, card_id, pitch_id, user_id, order_date, status,end_order_date) 
							VALUES
						($1,$2,$3,$4,$5,$6,to_timestamp($7) at time zone 'GMT',$8,to_timestamp($9) at time zone 'GMT') RETURNING id`, orderTable)

	err = tx.QueryRowx(query, order.UserName, order.PhoneNumber, order.ExtraInfo, order.CardId, order.PitchId, order.UserId, order.OrderDate, order.Status, total+int(order.OrderDate)).Scan(&id)

	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			return 0, fmt.Errorf("repository.Create:tx1 %w", txErr)
		}
		return 0, fmt.Errorf("repository.Create:err1 %w", err)
	}

	for _, val := range order.ServiceIds {
		queryService := fmt.Sprintf(
			`INSERT INTO
							%s
						(service_id,order_id) 
							VALUES
						($1,$2)`, orderServiceTable)

		_, err = tx.Exec(queryService, val, id)

		if err != nil {
			txErr := tx.Rollback()
			if txErr != nil {
				return 0, fmt.Errorf("repository.Create:tx2 %w", txErr)
			}
			return 0, fmt.Errorf("repository.Create:err2 %w", err)
		}

	}
	for _, val := range order.Times {
		queryService := fmt.Sprintf(
			`INSERT INTO
							%s
						(order_work_time,order_id) 
							VALUES
						($1,$2)`, orderTimeTable)

		_, err = tx.Exec(queryService, val, id)
		if err != nil {
			txErr := tx.Rollback()
			if txErr != nil {
				return 0, fmt.Errorf("repository.Create:tx3 %w", txErr)
			}
			return 0, fmt.Errorf("repository.Create:err3 %w", err)
		}

	}

	txErr := tx.Commit()
	if txErr != nil {
		return 0, fmt.Errorf("repository.Create: %w", txErr)
	}
	return id, nil
}

func (o *OrderRepos) GetAll(ctx *fiber.Ctx, page domain.Pagination, info domain.UserInfo, order domain.FilterForOrder) (*domain.GetAllResponses, error) {

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

	if order.OrderDate != 0 {
		forCheckValues = append(forCheckValues, fmt.Sprintf("o.order_date = to_timestamp(%f) at time zone 'GMT'", order.OrderDate))
	}

	switch order.OrderStatus {
	case 1:
		forCheckValues = append(forCheckValues, fmt.Sprintf("o.end_order_date < to_timestamp(%d) at time zone 'GMT'", time.Now().Unix()))
	case 2:
		forCheckValues = append(forCheckValues, fmt.Sprintf("o.end_order_date >= to_timestamp(%d) at time zone 'GMT'", time.Now().Unix()))
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
					o.status,
					o.first_name,
					o.phone_number,
					o.card_id,
					p.price,
					p.pitch_type,
					p.pitch_image,
					b.building_name,
					b.address
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
				%s ORDER BY o.id ASC LIMIT $1 OFFSET $2`, orderTable, pitchTable, buildingTable, setValues)

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

		queryTimes := fmt.Sprintf("SELECT id,order_work_time FROM %s WHERE order_id = $1 ", orderTimeTable)

		err = o.db.Select(&value.TimeInput, queryTimes, value.Id)

		if err != nil {
			return nil, fmt.Errorf("repository.GetAll: %w", err)

		}

		queryServices := fmt.Sprintf(
			`SELECT
						price,service_name
					FROM
						%s os
					INNER JOIN
						%s ser
					ON 
						os.service_id = ser.id
					WHERE
						order_id = $1 `, orderServiceTable, serviceTable)

		err = o.db.Select(&value.ServiceInput, queryServices, value.Id)

		if err != nil {
			return nil, fmt.Errorf("repository.GetAll: %w", err)

		}

	}

	ans := domain.GetAllResponses{
		Data:     inp,
		PageInfo: pages,
	}
	return &ans, nil
}

func (o *OrderRepos) GetAllBookTime(ctx *fiber.Ctx, times domain.FilterForOrderTimes) (*domain.GetAllResponses, error) {
	var (
		forCheckValues []string
		whereClause    string
		setValues      string
		caseValue      string
	)

	_, cancel := context.WithTimeout(ctx.Context(), 4*time.Second)

	defer cancel()

	if times.OrderDate != 0 {
		caseValue = fmt.Sprintf(
			`,CASE
						WHEN ot.id != 0 THEN true
						ELSE false
					END is_booked`)

		forCheckValues = append(forCheckValues, fmt.Sprintf("o.order_date = to_timestamp(%f) at time zone 'GMT'", times.OrderDate))

	}
	if times.BuildingId != 0 {
		forCheckValues = append(forCheckValues, fmt.Sprintf("t.building_id = %d", times.BuildingId))
	}

	whereClause = strings.Join(forCheckValues, " AND ")

	if len(forCheckValues) != 0 {
		setValues = " WHERE " + whereClause
	}

	var inp []*domain.OrderTime

	query := fmt.Sprintf(
		`select
					t.work_time
					%s
					from
						%s t
							LEFT JOIN
							%s b
								on
									b.id = t.building_id
							LEFT JOIN
							%s p
								on
									b.id = p.building_id
							LEFT JOIN 
								%s o
								on
									p.id = o.pitch_id
							LEFT JOIN 
								%s ot 
							on 
								(o.id = ot.order_id and t.work_time = ot.order_work_time)
							%s
					ORDER BY o.id ;`, caseValue, timeTable, buildingTable, pitchTable, orderTable, orderTimeTable, setValues)

	err := o.db.Select(&inp, query)

	if err != nil {
		return nil, fmt.Errorf("repository.GetAllBookTime: %w", err)
	}
	ans := domain.GetAllResponses{
		Data: inp,
	}
	return &ans, nil
}
