package targetdb

import (
	"context"
	"database/sql"

	"github.com/sirupsen/logrus"
	_ "golang.org/x/net/ipv4"

	"github.com/mojodojo101/c2server/pkg/models"
	"github.com/mojodojo101/c2server/pkg/target"
)

type sqlTargetRepo struct {
	DB *sql.DB
}

//BEGIN struct exported methods

func NewSQLRepo(db *sql.DB) target.Repository {
	return &sqlTargetRepo{
		DB: db,
	}
}

func (r *sqlTargetRepo) CreateTable(ctx context.Context) error {
	query := `CREATE TABLE IF NOT EXISTS target (
		id BIGSERIAL PRIMARY KEY,
		ipv4 CHAR(16),
		ipv6 CHAR(39),
		host_name CHAR(128),
		path CHAR(128),
		cmd_id BIGINT,
		created_at TIMESTAMP,
		updated_at TIMESTAMP)`
	_, err := r.execQuery(ctx, query)

	return err
}

func (r *sqlTargetRepo) DropTable(ctx context.Context) error {
	query := `DROP TABLE target`
	_, err := r.execQuery(ctx, query)

	return err
}

func (r *sqlTargetRepo) GetByID(ctx context.Context, id int64) (*models.Target, error) {
	query := `SELECT * FROM target WHERE id=$1`
	return r.getOneItem(ctx, query, id)
}

func (r *sqlTargetRepo) DeleteByID(ctx context.Context, id int64) error {
	query := `DELETE FROM target WHERE id=$1`
	_, err := r.execQuery(ctx, query, id)
	return err
}

func (r *sqlTargetRepo) Update(ctx context.Context, t *models.Target) error {
	query := `UPDATE target SET ipv4=$2,ipv6=$3,host_name=$4,path=$5,cmd_id=$6,updated_at=$7 WHERE id=$1`

	_, err := r.execQuery(ctx, query, t.Id, t.Ipv4, t.Ipv6, t.HostName, t.Path, t.CmdId, t.UpdatedAt)

	return err
}

func (r *sqlTargetRepo) CreateNewTarget(ctx context.Context, t *models.Target) error {

	query := `INSERT INTO target VALUES ( nextval('target_id_seq') ,$1,$2,$3,$4,$5,$6,$7) returning id`

	key, err := r.addItem(
		ctx,
		query,
		&t.Ipv4,
		&t.Ipv6,
		&t.HostName,
		&t.Path,
		&t.CmdId,
		&t.CreatedAt,
		&t.UpdatedAt)

	if err != nil {
		logrus.Error(err)
		return err
	}

	t.Id = key

	return nil

}

//END struct exported methods

func (r *sqlTargetRepo) getOneItem(ctx context.Context, query string, args ...interface{}) (*models.Target, error) {
	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	row := stmt.QueryRowContext(ctx, args...)

	t := &models.Target{}

	err = row.Scan(
		&t.Id,
		&t.Ipv4,
		&t.Ipv6,
		&t.HostName,
		&t.Path,
		&t.CmdId,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return t, nil
}

func (r *sqlTargetRepo) addItem(ctx context.Context, query string, args ...interface{}) (int64, error) {
	tx, err := r.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		logrus.Error(err)
		return 0, err
	}
	var id int64

	err = tx.QueryRowContext(ctx, query, args...).Scan(&id)

	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			logrus.Error(err)
		}
		logrus.Error(err)
	}

	if err := tx.Commit(); err != nil {
		logrus.Error(err)
	}

	return id, nil

}

// executes command and returns rows affected and error
func (r *sqlTargetRepo) execQuery(ctx context.Context, query string, args ...interface{}) (int64, error) {

	tx, err := r.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		logrus.Error(err)
		return 0, err
	}

	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			logrus.Error(err)
		}
		logrus.Error(err)
	}

	if err := tx.Commit(); err != nil {
		logrus.Error(err)
	}
	rows, err := result.RowsAffected()

	return rows, nil

}
