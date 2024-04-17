package postgresql

import (
	"github/user_service_evrone_microservces/internal/entity"
	"github/user_service_evrone_microservces/internal/pkg/postgres"
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
)

const (
	usersTableName      = "users"
	usersServiceName    = "usersService"
	usersSpanRepoPrefix = "usersRepo"
)

type usersRepo struct {
	tableName string
	db        *postgres.PostgresDB
}

func NewUsersRepo(db *postgres.PostgresDB) *usersRepo {
	return &usersRepo{
		tableName: usersTableName,
		db:        db,
	}
}

func (p *usersRepo) usersSelectQueryPrefix() squirrel.SelectBuilder {
	return p.db.Sq.Builder.
		Select(
			"guid",
			"full_name",
			"username",
			"email",
			"password",
			"bio",
			"website",
			"role",
			"created_at",
			"updated_at",
		).From(p.tableName)
}

func (p usersRepo) Create(ctx context.Context, users *entity.Users) error {
	// ctx, span := otlp.Start(ctx, usersServiceName, usersSpanRepoPrefix+"Create")
	// defer span.End()

	data := map[string]any{
		"guid":         users.GUID,
		"full_name":      users.FullName,
		"username":      users.Username,
		"email":        users.Email,
		"password":         users.Password,
		"bio":        users.Bio,
		"website":       users.Website,
		"role":  users.Role,
		"created_at":   users.CreatedAt,
		"updated_at":   users.UpdatedAt,
	}
	query, args, err := p.db.Sq.Builder.Insert(p.tableName).SetMap(data).ToSql()
	if err != nil {
		return p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "create"))
	}

	_, err = p.db.Exec(ctx, query, args...)
	if err != nil {
		return p.db.Error(err)
	}

	return nil
}

func (p usersRepo) Get(ctx context.Context, params map[string]string) (*entity.Users, error) {
	// ctx, span := otlp.Start(ctx, usersServiceName, usersSpanRepoPrefix+"Get")
	// defer span.End()

	var (
		users entity.Users
	)

	queryBuilder := p.usersSelectQueryPrefix()

	for key, value := range params {
		if key == "guid" {
			queryBuilder = queryBuilder.Where(p.db.Sq.Equal(key, value))
		}
	}
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "get"))
	}
	if err = p.db.QueryRow(ctx, query, args...).Scan(
		&users.GUID,
		&users.FullName,
		&users.Username,
		&users.Email,
		&users.Password,
		&users.Bio,
		&users.Website,
		&users.Role,
		&users.CreatedAt,
		&users.UpdatedAt,
	); err != nil {
		return nil, p.db.Error(err)
	}

	return &users, nil
}

func (p usersRepo) List(ctx context.Context, limit, offset uint64, filter map[string]string) ([]*entity.Users, error) {
	// ctx, span := otlp.Start(ctx, usersServiceName, usersSpanRepoPrefix+"List")
	// defer span.End()

	var (
		users []*entity.Users
	)
	queryBuilder := p.usersSelectQueryPrefix()

	if limit != 0 {
		queryBuilder = queryBuilder.Limit(limit).Offset(offset)
	}
	fmt.Println(filter)
	for key, value := range filter {
		if key == "type_id" || key == "lang" || key == "status" {
			queryBuilder = queryBuilder.Where(p.db.Sq.Equal(key, value))
			continue
		}
		if key == "created_at" {
			queryBuilder = queryBuilder.Where("created_at=?", value)
			continue
		}
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "list"))
	}

	rows, err := p.db.Query(ctx, query, args...)
	if err != nil {
		return nil, p.db.Error(err)
	}
	defer rows.Close()
	users = make([]*entity.Users, 0)
	for rows.Next() {
		var user entity.Users
		if err = rows.Scan(
			&user.GUID,
			&user.FullName,
			&user.Username,
			&user.Email,
			&user.Password,
			&user.Bio,
			&user.Website,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, p.db.Error(err)
		}
		users = append(users, &user)
	}

	return users, nil
}

func (p usersRepo) Delete(ctx context.Context, guid string) error {
	// ctx, span := otlp.Start(ctx, usersServiceName, usersSpanRepoPrefix+"Delete")
	// defer span.End()

	sqlStr, args, err := p.db.Sq.Builder.
		Delete(p.tableName).
		Where(p.db.Sq.Equal("guid", guid)).
		ToSql()
	if err != nil {
		return p.db.ErrSQLBuild(err, p.tableName+" delete")
	}

	commandTag, err := p.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return p.db.Error(err)
	}

	if commandTag.RowsAffected() == 0 {
		return p.db.Error(fmt.Errorf("no sql rows"))
	}

	return nil
}