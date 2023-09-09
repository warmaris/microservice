package looncan

import (
	"context"
	"fmt"
	"microservice/internal/database"
	"strings"
	"time"
)

type mysqlStorage struct {
	conn database.Connection
}

func NewStorage(conn database.Connection) Storage {
	return &mysqlStorage{
		conn: conn,
	}
}

func (m *mysqlStorage) getAllForParent(ctx context.Context, parentID uint64, parentType ParentType) ([]Looncan, error) {
	const query = `SELECT id, name, value, parent_id, parent_type, created_at, updated_at FROM looncan WHERE parent_id = ? AND parent_type = ?`

	var rows []model
	err := m.conn.Select(ctx, &rows, query, parentID, parentType)

	if err != nil {
		return nil, err
	}

	res := make([]Looncan, 0, len(rows))
	for _, row := range rows {
		var updatedAt *time.Time
		if row.UpdatedAt.Valid {
			updatedAt = new(time.Time)
			*updatedAt = row.UpdatedAt.Time
		}
		res = append(res, Looncan{
			ID:         row.ID,
			Name:       row.Name,
			Value:      row.Value,
			ParentID:   row.ParentID,
			ParentType: ParentType(row.ParentType),
			CreatedAt:  row.CreatedAt,
			UpdatedAt:  updatedAt,
		})
	}

	return res, nil
}

func (m *mysqlStorage) list(ctx context.Context) ([]Looncan, error) {
	const query = `SELECT id, name, value, parent_id, parent_type, created_at, updated_at FROM looncan ORDER BY id`

	var rows []model
	err := m.conn.Select(ctx, &rows, query)

	if err != nil {
		return nil, err
	}

	res := make([]Looncan, 0, len(rows))
	for _, row := range rows {
		var updatedAt *time.Time
		if row.UpdatedAt.Valid {
			updatedAt = new(time.Time)
			*updatedAt = row.UpdatedAt.Time
		}
		res = append(res, Looncan{
			ID:         row.ID,
			Name:       row.Name,
			Value:      row.Value,
			ParentID:   row.ParentID,
			ParentType: ParentType(row.ParentType),
			CreatedAt:  row.CreatedAt,
			UpdatedAt:  updatedAt,
		})
	}

	return res, nil
}

func (m *mysqlStorage) create(ctx context.Context, entities []Looncan) error {
	const query = `INSERT INTO looncan (name, value, parent_id, parent_type) VALUES %s`

	placeholders := make([]string, len(entities))
	args := make([]any, 0, 4*len(entities))

	for i, l := range entities {
		placeholders[i] = "(?, ?, ?, ?)"
		args = append(args, l.Name, l.Value, l.ParentID, l.ParentType)
	}

	return m.conn.Exec(ctx, fmt.Sprintf(query, strings.Join(placeholders, ", ")), args...)
}
