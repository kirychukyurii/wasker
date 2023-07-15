package repository

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	scan "github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/kirychukyurii/wasker/internal/app/directory/model"
	"github.com/kirychukyurii/wasker/internal/lib"
	"github.com/kirychukyurii/wasker/internal/pkg/server/interceptor/requestid"

	"github.com/kirychukyurii/wasker/internal/errors"
	"github.com/kirychukyurii/wasker/internal/pkg/db"
	"github.com/kirychukyurii/wasker/internal/pkg/log"
)

type AuthRepository struct {
	db     db.Database
	logger log.Logger
}

func NewAuthRepository(db db.Database, logger log.Logger) AuthRepository {
	return AuthRepository{
		db:     db,
		logger: logger,
	}
}

func (a AuthRepository) Login(ctx context.Context, login *model.UserSession) error {
	q := a.db.Dialect().Insert("auth_user_session").
		Columns("user_id", "network_ip", "access_token", "created_at", "expires_at").
		Values(login.User.Id, login.NetworkIp, login.AccessToken, login.CreatedAt, login.ExpiresAt)

	sql, args, err := q.ToSql()
	if err != nil {
		return errors.NewInternalError(errors.AppError{
			Message: errors.ErrDatabaseInternalError.Error(),
			Details: errors.AppErrorDetail{
				ErrId:     "repository.auth.login.build_query",
				Err:       err,
				RequestId: lib.FromContext(ctx, requestid.XRequestIDCtxKey{}).(string),
			},
		})
	}

	if _, err = a.db.Pool.Exec(ctx, sql, args...); err != nil {
		return errors.NewInternalError(errors.AppError{
			Message: errors.ErrDatabaseInternalError.Error(),
			Details: errors.AppErrorDetail{
				ErrId:     "repository.auth.login.build_query",
				Err:       err,
				RequestId: lib.FromContext(ctx, requestid.XRequestIDCtxKey{}).(string),
			},
		})
	}

	return nil
}

func (a AuthRepository) Authn(ctx context.Context, token string) (*model.UserSession, error) {
	var session model.UserSession

	q := a.db.Dialect().Select("id", `user_id "user.id"`, "created_at", "access_token", "network_ip", "expires_at").
		From("auth_user_session").
		Where(sq.Eq{"access_token": token})

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, errors.NewInternalError(errors.AppError{
			Message: errors.ErrDatabaseInternalError.Error(),
			Details: errors.AppErrorDetail{
				ErrId:     "repository.auth.authn.build_query",
				Err:       err,
				RequestId: lib.FromContext(ctx, requestid.XRequestIDCtxKey{}).(string),
			},
		})
	}

	if err := scan.Get(ctx, a.db.Pool, &session, sql, args...); err != nil {
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
				ErrId:     "repository.auth.authn.exec_query",
				Err:       err,
				RequestId: lib.FromContext(ctx, requestid.XRequestIDCtxKey{}).(string),
			},
		})
	}

	return &session, nil
}

func (a AuthRepository) Authz(ctx context.Context, userId int64, service, method string) (endpoint, permission uint8, err error) {
	fullMethod := fmt.Sprintf("%s%s", service, method)

	q := a.db.Dialect().Select("se.bit as endpoint_bit", "rp.bit as permission_bit").From("auth_role_permission rp").
		InnerJoin("auth_user u on u.role_id = rp.role_id").
		InnerJoin("auth_scope_endpoint se on se.scope_id = rp.scope_id").
		Where(sq.Eq{"u.id": userId}).
		Where(sq.Eq{"se.name": fullMethod})

	sql, args, err := q.ToSql()
	if err != nil {
		a.logger.FromContext(ctx).Log.Error().Err(err).Msg(errors.ErrDatabaseBuildSql.Error())

		return 0, 0, errors.ErrDatabaseInternalError
	}

	if err := a.db.Pool.QueryRow(ctx, sql, args...).Scan(&endpoint, &permission); err != nil {
		a.logger.FromContext(ctx).Log.Error().Err(err).Msg(errors.ErrDatabaseQueryRow.Error())

		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return 0, 0, errors.ErrDatabaseRecordNotFound
		default:
			return 0, 0, errors.ErrDatabaseInternalError
		}
	}

	return endpoint, permission, nil
}
