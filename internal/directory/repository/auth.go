package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/kirychukyurii/wasker/internal/directory/model"

	sq "github.com/Masterminds/squirrel"
	scan "github.com/georgysavva/scany/v2/pgxscan"

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
		a.logger.FromContext(ctx).Log.Error().Err(err).Msg(errors.ErrDatabaseBuildSql.Error())

		return errors.ErrDatabaseInternalError
	}

	if _, err = a.db.Pool.Exec(ctx, sql, args...); err != nil {
		a.logger.FromContext(ctx).Log.Error().Err(err).Msg(errors.ErrDatabaseQueryRow.Error())

		return errors.ErrDatabaseInternalError
	}

	return nil
}

func (a AuthRepository) VerifyToken(ctx context.Context, token string) (*model.UserSession, error) {
	var session model.UserSession

	q := a.db.Dialect().Select("id", `user_id "user.id"`, "created_at", "access_token", "network_ip", "expires_at").
		From("auth_user_session").
		Where(sq.Eq{"access_token": token})

	sql, args, err := q.ToSql()
	if err != nil {
		a.logger.FromContext(ctx).Log.Error().Err(err).Msg(errors.ErrDatabaseBuildSql.Error())

		return nil, errors.ErrDatabaseInternalError
	}

	if err := scan.Get(ctx, a.db.Pool, &session, sql, args...); err != nil {
		a.logger.FromContext(ctx).Log.Error().Err(err).Msg(errors.ErrDatabaseQueryRow.Error())

		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, errors.ErrDatabaseRecordNotFound
		default:
			return nil, errors.ErrDatabaseInternalError
		}
	}

	return &session, nil
}

func (a AuthRepository) VerifyPermission(ctx context.Context, userId uint64, service, method string) (endpoint, permission uint8, err error) {
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
