package v1alpha1

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/pkg/errors"

	"github.com/kirychukyurii/wasker/internal/model"
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
	u := model.User{}

	q := a.db.Dialect().Select("u.id", "coalesce(u.name, u.user_name) AS name", "u.user_name", "u.password", "u.email",
		"u.created_at", `u.created_by "created_by.id"`, `coalesce(c.name, c.user_name) "created_by.name"`,
		"u.updated_at", `u.updated_by "updated_by.id"`, `coalesce(upd.name, upd.user_name) "updated_by.name"`,
		`u.role_id "role.id"`, `r.name "role.name"`).From("auth_user u").
		LeftJoin("auth_role r ON r.id = u.role_id").
		InnerJoin("auth_user c ON c.id = u.created_by").
		InnerJoin("auth_user upd ON upd.id = u.updated_by").
		Where(sq.Eq{"u.deleted_at": nil}).Where(sq.Eq{"u.id": userId})

	/*
			if roleId != 0 {
				sqlOwner := sq.Select("distinct owner.id").PlaceholderFormat(sq.Dollar).
					From("auth_user owner").
					LeftJoin("auth_map_acl ama ON ama.grantor_role_id = owner.role_id").
					Where(sq.Eq{"ama.grantee_role_id": roleId})
				sql = sql.Where(sq.Or{sq.Eq{"u.created_by": userId}, sqlOwner.Prefix("u.created_by IN (").Suffix(")")})
			} else {
		q = q.Where(sq.Eq{"u.created_by": userId})
		} */

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "building SQL statement")
	}

	if err := pgxscan.Get(ctx, a.db.Pool, &u, sql, args...); err != nil {
		return nil, errors.Wrap(err, "querying rows")
	}

	/*
		row, err := a.db.Pool.Query(ctx, sql, args...)
		if err != nil {
			return nil, errors.Wrap(err, "failed querying rows")
		}

		var user model.User
		if err := pgxscan.NewScanner(row).Scan(&user); err != nil {
			return nil, errors.Wrap(err, "failed collecting rows")
		}



		user, err := pgx.CollectOneRow(row, pgx.RowToStructByName[model.User])
		if err != nil {
			return nil, errors.Wrap(err, "failed collecting rows")
		}



		Scan(&u.Id, &u.Name, &u.UserName, &u.Password, &u.Email, &u.CreatedAt, &u.CreatedBy.Id, &u.CreatedBy.Name,
			&u.UpdatedAt, &u.UpdatedBy.Id, &u.UpdatedBy.Name, &u.Role.Id, &u.Role.Name)
	*/

	return &u, nil
}
