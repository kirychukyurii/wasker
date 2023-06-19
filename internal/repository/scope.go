package repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/kirychukyurii/wasker/internal/errors"
	"github.com/kirychukyurii/wasker/internal/model"
	"github.com/kirychukyurii/wasker/internal/model/dto"
	"github.com/kirychukyurii/wasker/internal/pkg/db"
	"github.com/kirychukyurii/wasker/internal/pkg/log"
)

type ScopeRepository struct {
	db     db.Database
	logger log.Logger
}

func NewScopeRepository(db db.Database, logger log.Logger) ScopeRepository {
	return ScopeRepository{
		db:     db,
		logger: logger,
	}
}

func (a ScopeRepository) Create(ctx context.Context, scope *model.Scope) (uint64, error) {
	var rowId uint64

	q := a.db.Dialect().Insert("auth_scope").
		Columns("name").Values(scope.Name).Suffix("RETURNING id")

	sql, args, err := q.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "build SQL statement")
	}

	err = a.db.Pool.QueryRow(ctx, sql, args...).Scan(&rowId)
	if err != nil {
		return 0, err
	}

	return rowId, nil
}

func (a ScopeRepository) Query(ctx context.Context, param *model.ScopeQueryParam) (*model.ScopeQueryResult, error) {
	var list model.Scopes
	var scope model.Scope
	var pagination dto.Pagination

	q := a.db.Dialect().Select("id", "name").
		From("auth_scope")

	if v := param.Query.Id; v != 0 {
		q = q.Where(sq.Eq{"id": v})
	}

	if v := param.Query.Name; v != "" {
		q = q.Where(sq.Eq{"name": v})
	}

	if v := param.Query.Query; v != "" {
		q = q.Where(sq.Or{
			sq.Like{"name": v},
		})
	}

	q = q.OrderBy(fmt.Sprintf("%s", param.Order.ParseOrder()))
	current, pageSize := param.Pagination.GetCurrent(), param.Pagination.GetPageSize()
	if current > 0 && pageSize > 0 {
		offset := (current - 1) * pageSize
		q = q.Offset(offset).Limit(pageSize)
	} else if pageSize > 0 {
		q = q.Limit(pageSize)
	}

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "build SQL statement")
	}

	rows, err := a.db.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "query rows")
	}

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&scope.Id, &scope.Name)
		if err != nil {
			return nil, errors.Wrap(err, "scan rows")
		}

		list = append(list, &scope)
	}
	if rows.Err() != nil {
		return nil, errors.Wrap(rows.Err(), "read rows")
	}

	pagination.Current = current
	pagination.PageSize = pageSize

	qr := &model.ScopeQueryResult{
		Pagination: &pagination,
		List:       list,
	}

	return qr, nil
}

func (a ScopeRepository) CreateEndpoint(ctx context.Context, endpoint *model.ScopeEndpoint) (uint64, error) {
	var rowId uint64

	q := a.db.Dialect().Insert("auth_scope_endpoint").
		Columns("name", "bit", "scope_id").Values(endpoint.Name, endpoint.Bit, endpoint.Scope.Id).
		Suffix("RETURNING id")

	sql, args, err := q.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "build SQL statement")
	}

	err = a.db.Pool.QueryRow(ctx, sql, args...).Scan(&rowId)
	if err != nil {
		return 0, err
	}

	return rowId, nil
}

func (a ScopeRepository) QueryEndpoint(ctx context.Context, param *model.ScopeEndpointQueryParam) (*model.ScopeEndpointQueryResult, error) {
	var list model.ScopeEndpoints
	var endpoint model.ScopeEndpoint
	var pagination dto.Pagination

	q := a.db.Dialect().Select("se.id", "se.name", "se.bit", `s.id as "scope.id"`, `s.name as "scope.name"`).
		From("auth_scope_endpoint se").
		InnerJoin("auth_scope s on s.id = se.scope_id")

	if v := param.Query.Id; v != 0 {
		q = q.Where(sq.Eq{"se.id": v})
	}

	if v := param.Query.Name; v != "" {
		q = q.Where(sq.Eq{"se.name": v})
	}

	if v := param.Query.Query; v != "" {
		q = q.Where(sq.Or{
			sq.Like{"se.name": v},
		})
	}

	q = q.OrderBy(fmt.Sprintf("se.%s", param.Order.ParseOrder()))
	current, pageSize := param.Pagination.GetCurrent(), param.Pagination.GetPageSize()
	if current > 0 && pageSize > 0 {
		offset := (current - 1) * pageSize
		q = q.Offset(offset).Limit(pageSize)
	} else if pageSize > 0 {
		q = q.Limit(pageSize)
	}

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "build SQL statement")
	}

	rows, err := a.db.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "query rows")
	}

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&endpoint.Id, &endpoint.Name, &endpoint.Bit, &endpoint.Scope.Id, &endpoint.Scope.Name)
		if err != nil {
			return nil, errors.Wrap(err, "scan rows")
		}

		list = append(list, &endpoint)
	}
	if rows.Err() != nil {
		return nil, errors.Wrap(rows.Err(), "read rows")
	}

	pagination.Current = current
	pagination.PageSize = pageSize

	qr := &model.ScopeEndpointQueryResult{
		Pagination: &pagination,
		List:       list,
	}

	return qr, nil
}
