package repository

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/kirychukyurii/wasker/internal/model"

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

func (a AuthRepository) CheckToken(ctx context.Context, token string) (*model.UserSession, error) {
	var s model.UserSession

	q := a.db.Dialect().Select("user_id").From("auth_user_session").
		Where(sq.Eq{"access_token": token})

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "build SQL statement")
	}

	if err := a.db.Pool.QueryRow(ctx, sql, args...).Scan(&s.Id, &s.CreatedAt, s.AccessToken, s.NetworkIp, &s.ExpiresAt); err != nil {
		return nil, errors.Wrap(err, "query rows")
	}

	return &s, nil
}
