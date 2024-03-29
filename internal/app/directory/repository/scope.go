package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	scan "github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"

	"github.com/kirychukyurii/wasker/internal/app/directory/model"
	"github.com/kirychukyurii/wasker/internal/errors"
	"github.com/kirychukyurii/wasker/internal/lib"
	gmodel "github.com/kirychukyurii/wasker/internal/model"
	"github.com/kirychukyurii/wasker/internal/pkg/db"
	"github.com/kirychukyurii/wasker/internal/pkg/log"
	"github.com/kirychukyurii/wasker/internal/pkg/server/interceptor/requestid"
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

func (a ScopeRepository) Create(ctx context.Context, scope *model.Scope) (int64, error) {
	var rowId int64

	q := a.db.Dialect().Insert("auth_scope").
		Columns("name").Values(scope.Name).Suffix("RETURNING id")

	sql, args, err := q.ToSql()
	if err != nil {
		return 0, errors.NewInternalError(errors.AppError{
			Message: errors.ErrDatabaseInternalError.Error(),
			Details: errors.AppErrorDetail{
				Err:       err,
				ErrReason: errors.ErrBuildQueryReason,
				ErrDomain: "repository.scope.create",
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
				ErrDomain: "repository.scope.create",
				RequestId: lib.FromContext(ctx, requestid.XRequestIDCtxKey{}).(string),
			},
		})
	}

	return rowId, nil
}

func (a ScopeRepository) Query(ctx context.Context, param *model.ScopeQueryParam) (*model.ScopeQueryResult, error) {
	var list model.Scopes
	var pagination gmodel.Pagination

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

	q = q.OrderBy(param.Order.Parse())
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
				ErrDomain: "repository.scope.query",
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
				ErrDomain: "repository.scope.query",
				RequestId: lib.FromContext(ctx, requestid.XRequestIDCtxKey{}).(string),
			},
		})
	}

	pagination.Current = current
	pagination.PageSize = pageSize

	qr := &model.ScopeQueryResult{
		Pagination: &pagination,
		List:       list,
	}

	return qr, nil
}
