package repository

import (
	"context"
	"fmt"
	"github.com/kirychukyurii/wasker/internal/lib"
	"github.com/kirychukyurii/wasker/internal/pkg/server/interceptor/requestid"

	sq "github.com/Masterminds/squirrel"
	scan "github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"

	"github.com/kirychukyurii/wasker/internal/app/directory/model"
	"github.com/kirychukyurii/wasker/internal/errors"
	gmodel "github.com/kirychukyurii/wasker/internal/model"
)

func (a ScopeRepository) CreateEndpoint(ctx context.Context, endpoint *model.ScopeEndpoint) (uint64, error) {
	var rowId uint64

	q := a.db.Dialect().Insert("auth_scope_endpoint").
		Columns("name", "bit", "scope_id").Values(endpoint.Name, endpoint.Bit, endpoint.Scope.Id).
		Suffix("RETURNING id")

	sql, args, err := q.ToSql()
	if err != nil {
		return 0, errors.NewInternalError(errors.AppError{
			Message: errors.ErrDatabaseInternalError.Error(),
			Details: errors.AppErrorDetail{
				Err:       err,
				ErrReason: errors.ErrBuildQueryReason,
				ErrDomain: "repository.scope_endpoint.create",
				RequestId: lib.FromContext(ctx, requestid.XRequestIDCtxKey{}).(string),
			},
		})
	}

	err = a.db.Pool.QueryRow(ctx, sql, args...).Scan(&rowId)
	if err != nil {
		return 0, errors.NewInternalError(errors.AppError{
			Message: errors.ErrDatabaseInternalError.Error(),
			Details: errors.AppErrorDetail{
				Err:       err,
				ErrReason: errors.ErrExecQueryReason,
				ErrDomain: "repository.scope_endpoint.create",
				RequestId: lib.FromContext(ctx, requestid.XRequestIDCtxKey{}).(string),
			},
		})
	}

	return rowId, nil
}

func (a ScopeRepository) QueryEndpoint(ctx context.Context, param *model.ScopeEndpointQueryParam) (*model.ScopeEndpointQueryResult, error) {
	var list model.ScopeEndpoints
	var pagination gmodel.Pagination

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

	q = q.OrderBy(fmt.Sprintf("se.%s", param.Order.Parse()))
	current, pageSize := param.Pagination.GetCurrent(), param.Pagination.GetPageSize()
	if current > 0 && pageSize > 0 {
		offset := (current - 1) * pageSize
		q = q.Offset(offset).Limit(pageSize)
	} else if pageSize > 0 {
		q = q.Limit(pageSize)
	}

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, errors.NewInternalError(errors.AppError{
			Message: errors.ErrDatabaseInternalError.Error(),
			Details: errors.AppErrorDetail{
				Err:       err,
				ErrReason: errors.ErrBuildQueryReason,
				ErrDomain: "repository.scope_endpoint.query",
				RequestId: lib.FromContext(ctx, requestid.XRequestIDCtxKey{}).(string),
			},
		})
	}

	if err = scan.Select(ctx, a.db.Pool, &list, sql, args...); err != nil {
		var dbErr error
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			dbErr = errors.ErrDatabaseRecordNotFound
		default:
			dbErr = errors.ErrDatabaseInternalError
		}

		return nil, errors.NewInternalError(errors.AppError{
			Message: dbErr.Error(),
			Details: errors.AppErrorDetail{
				Err:       err,
				ErrReason: errors.ErrExecQueryReason,
				ErrDomain: "repository.scope_endpoint.query",
				RequestId: lib.FromContext(ctx, requestid.XRequestIDCtxKey{}).(string),
			},
		})
	}

	pagination.Current = current
	pagination.PageSize = pageSize

	qr := &model.ScopeEndpointQueryResult{
		Pagination: &pagination,
		List:       list,
	}

	return qr, nil
}
