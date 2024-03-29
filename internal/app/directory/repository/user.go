package repository

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	scan "github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/kirychukyurii/wasker/internal/app/directory/model"
	"github.com/kirychukyurii/wasker/internal/lib"
	model2 "github.com/kirychukyurii/wasker/internal/model"
	"github.com/kirychukyurii/wasker/internal/pkg/server/interceptor/requestid"

	"github.com/kirychukyurii/wasker/internal/errors"
	"github.com/kirychukyurii/wasker/internal/pkg/db"
	"github.com/kirychukyurii/wasker/internal/pkg/log"
)

type UserRepository struct {
	db     db.Database
	logger log.Logger
}

func NewUserRepository(db db.Database, logger log.Logger) UserRepository {
	return UserRepository{
		db:     db,
		logger: logger,
	}
}

func (a UserRepository) ReadUser(ctx context.Context, userId int64) (*model.User, error) {
	var user model.User

	q := a.db.Dialect().Select("u.id", "coalesce(u.name, u.user_name) AS name", "u.user_name", "u.password", "u.email",
		"u.created_at", `u.created_by "created_by.id"`, `coalesce(c.name, c.user_name) "created_by.name"`,
		"u.updated_at", `u.updated_by "updated_by.id"`, `coalesce(upd.name, upd.user_name) "updated_by.name"`,
		`u.role_id "role.id"`, `r.name "role.name"`).From("auth_user u").
		LeftJoin("auth_role r ON r.id = u.role_id").
		InnerJoin("auth_user c ON c.id = u.created_by").
		InnerJoin("auth_user upd ON upd.id = u.updated_by").
		Where(sq.Eq{"u.deleted_at": nil}).Where(sq.Eq{"u.id": userId})

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, errors.NewInternalError(errors.AppError{
			Message: errors.ErrDatabaseInternalError.Error(),
			Details: errors.AppErrorDetail{
				Err:       err,
				ErrReason: errors.ErrBuildQueryReason,
				ErrDomain: "repository.user.read",
				RequestId: lib.FromContext(ctx, requestid.XRequestIDCtxKey{}).(string),
			},
		})
	}

	if err = scan.Get(ctx, a.db.Pool, &user, sql, args...); err != nil {
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
				ErrDomain: "repository.user.read",
				RequestId: lib.FromContext(ctx, requestid.XRequestIDCtxKey{}).(string),
			},
		})
	}

	return &user, nil
}

func (a UserRepository) Query(ctx context.Context, param *model.UserQueryParam) (*model.UserQueryResult, error) {
	var list model.Users
	var pagination model2.Pagination

	q := a.db.Dialect().Select("u.id", "coalesce(u.name, u.user_name) AS name", "u.user_name", "u.password", "u.email",
		"u.created_at", `u.created_by "created_by.id"`, `coalesce(c.name, c.user_name) "created_by.name"`,
		"u.updated_at", `u.updated_by "updated_by.id"`, `coalesce(upd.name, upd.user_name) "updated_by.name"`,
		`u.role_id "role.id"`, `r.name "role.name"`).From("auth_user u").
		LeftJoin("auth_role r ON r.id = u.role_id").
		InnerJoin("auth_user c ON c.id = u.created_by").
		InnerJoin("auth_user upd ON upd.id = u.updated_by").
		Where(sq.Eq{"u.deleted_at": nil})

	if v := param.Query.Id; v != 0 {
		q = q.Where(sq.Eq{"u.id": v})
	}

	if v := param.Query.Name; v != "" {
		q = q.Where(sq.Eq{"u.name": v})
	}

	if v := param.UserName; v != "" {
		q = q.Where(sq.Eq{"u.user_name": v})
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
				ErrDomain: "repository.user.query",
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
				ErrDomain: "repository.user.query",
				RequestId: lib.FromContext(ctx, requestid.XRequestIDCtxKey{}).(string),
			},
		})
	}

	pagination.Current = current
	pagination.PageSize = pageSize

	qr := &model.UserQueryResult{
		Pagination: &pagination,
		List:       list,
	}

	return qr, nil
}
