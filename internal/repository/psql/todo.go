package psql

import (
	"context"
	"fmt"
	"strconv"

	"github.com/abrahammegantoro/to-do-list-be/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TodoRepository struct {
	Conn *pgxpool.Pool
}

func NewTodoRepository(conn *pgxpool.Pool) *TodoRepository {
	return &TodoRepository{
		Conn: conn,
	}
}

func (t *TodoRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Todo, err error) {
	rows, err := t.Conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		td := domain.Todo{}
		err = rows.Scan(
			&td.ID,
			&td.Text,
			&td.Category,
			&td.Date,
			&td.PriorityLevel,
			&td.UserID,
			&td.Completed,
			&td.UpdatedAt,
			&td.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, td)
	}

	return
}

func (t *TodoRepository) Fetch(ctx context.Context, limit int64, offset int64) (res []domain.Todo, err error) {
	query := `SELECT id,text,category,date,priority_level,user_id,completed,updated_at,created_at FROM todos ORDER BY created_at DESC LIMIT $1 OFFSET $2`

	res, err = t.fetch(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	return
}

func (t *TodoRepository) GetByID(ctx context.Context, id int64) (res domain.Todo, err error) {
	query := `SELECT id,text,category,date,priority_level,user_id,completed,updated_at,created_at FROM todos WHERE ID = $1`

	list, err := t.fetch(ctx, query, id)
	if err != nil {
		return domain.Todo{}, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}

	return
}

func (t *TodoRepository) GetByUserID(ctx context.Context, userID int64, limit int64, offset int64, category *string, priorityLevel *string, keyword *string) (res []domain.Todo, err error) {
    query := `SELECT id, text, category, date, priority_level, user_id, completed, updated_at, created_at FROM todos WHERE user_id = $1`
    
    params := []interface{}{userID}
    paramIndex := 2

    if category != nil && *category != "" {
        query += ` AND category = $` + strconv.Itoa(paramIndex)
        params = append(params, *category)
        paramIndex++
    }
    if priorityLevel != nil && *priorityLevel != "" {
        query += ` AND priority_level = $` + strconv.Itoa(paramIndex)
        params = append(params, *priorityLevel)
        paramIndex++
    }
    if keyword != nil && *keyword != "" {
        query += ` AND text ILIKE '%' || $` + strconv.Itoa(paramIndex) + ` || '%'`
        params = append(params, *keyword)
        paramIndex++
    }

    query += ` ORDER BY created_at DESC LIMIT $` + strconv.Itoa(paramIndex) + ` OFFSET $` + strconv.Itoa(paramIndex+1)
    params = append(params, limit, offset)

    res, err = t.fetch(ctx, query, params...)
    if err != nil {
        return nil, err
    }

    return
}

func (t *TodoRepository) GetAllCategories(ctx context.Context) (res []string, err error) {
	query := `SELECT DISTINCT category FROM todos`

	rows, err := t.Conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var category string
		err = rows.Scan(&category)
		if err != nil {
			return nil, err
		}

		res = append(res, category)
	}

	return
}

func (t *TodoRepository) Store(ctx context.Context, td *domain.Todo) (err error) {
	query := `INSERT INTO todos (text, category, date, priority_level, user_id, updated_at, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7) returning id`

	err = t.Conn.QueryRow(ctx, query, td.Text, td.Category, td.Date, td.PriorityLevel, td.UserID, td.UpdatedAt, td.CreatedAt).Scan(&td.ID)
	if err != nil {
		return
	}

	return
}

func (t *TodoRepository) Delete(ctx context.Context, id int64) (err error) {
	query := `DELETE FROM todos WHERE id = $1`

	commandTag, err := t.Conn.Exec(ctx, query, id)
	if err != nil {
		return
	}

	rowsAfected := commandTag.RowsAffected()
	if rowsAfected != 1 {
		err = fmt.Errorf("weird  Behavior. Total Affected: %d", rowsAfected)
		return
	}

	return
}

func (t *TodoRepository) Update(ctx context.Context, td *domain.Todo) (err error) {
	query := `UPDATE todos SET text=$1, category=$2, date=$3, priority_level=$4, user_id=$5, completed=$6, updated_at=$7 WHERE id=$8`

	commandTag, err := t.Conn.Exec(ctx, query, td.Text, td.Category, td.Date, td.PriorityLevel, td.UserID, td.Completed, td.UpdatedAt, td.ID)
	if err != nil {
		return
	}

	rowsAfected := commandTag.RowsAffected()
	if rowsAfected != 1 {
		err = fmt.Errorf("weird  Behavior. Total Affected: %d", rowsAfected)
		return
	}

	return
}
