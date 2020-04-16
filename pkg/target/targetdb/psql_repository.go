package targetdb

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

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
func (r *sqlTargetRepo) GetByIpv4(ctx context.Context, ipv4 string) (*models.Target, error) {
	query := `SELECT * FROM target WHERE ipv4=$1`
	return r.getOneItem(ctx, query, ipv4)
}
func (r *sqlTargetRepo) GetAllTargets(ctx context.Context, amount int64) ([]models.Target, error) {
	query := `SELECT * FROM target`
	return r.getAllItems(ctx, query, amount)
}
func (r *sqlTargetRepo) DeleteByID(ctx context.Context, id int64) error {
	query := `DELETE FROM target WHERE id=$1`
	_, err := r.execQuery(ctx, query, id)
	return err
}

func (r *sqlTargetRepo) Update(ctx context.Context, t *models.Target) error {
	query := `UPDATE target SET ipv4=$2,ipv6=$3,host_name=$4,path=$5,updated_at=$6 WHERE id=$1`

	_, err := r.execQuery(ctx, query, t.Id, t.Ipv4, t.Ipv6, t.HostName, t.Path, t.UpdatedAt)

	return err
}

func (r *sqlTargetRepo) CreateNewTarget(ctx context.Context, t *models.Target) error {

	query := `INSERT INTO target VALUES ( nextval('target_id_seq') ,$1,$2,$3,$4,$5,$6) returning id`

	key, err := r.addItem(
		ctx,
		query,
		&t.Ipv4,
		&t.Ipv6,
		&t.HostName,
		&t.Path,
		&t.CreatedAt,
		&t.UpdatedAt)

	if err != nil {
		logrus.Error(err)
	}
	t.Id = key
	return nil

}

//END struct exported methods

func (r *sqlTargetRepo) getOneItem(ctx context.Context, query string, args ...interface{}) (*models.Target, error) {
	fmt.Printf("query for target = %v  with val %v", query, args)
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
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	t.Ipv4 = strings.TrimRight(t.Ipv4, " ")
	t.Ipv6 = strings.TrimRight(t.Ipv6, " ")
	t.Path = strings.TrimRight(t.Path, " ")
	t.HostName = strings.TrimRight(t.HostName, " ")

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
func (r *sqlTargetRepo) getAllItems(ctx context.Context, query string, amount int64, args ...interface{}) ([]models.Target, error) {
	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer rows.Close()
	ts := make([]models.Target, amount)
	for i, _ := range ts {
		rows.Next()
		err = rows.Scan(
			&ts[i].Id,
			&ts[i].Ipv4,
			&ts[i].Ipv6,
			&ts[i].HostName,
			&ts[i].Path,
			&ts[i].CreatedAt,
			&ts[i].UpdatedAt,
		)
		if err != nil {
			logrus.Error(err)
			if i != 0 {
				ts = append(ts[:i])
				return ts, nil
			}
			return nil, err
		}

		ts[i].HostName = strings.TrimRight(ts[i].HostName, " ")
		ts[i].Ipv4 = strings.TrimRight(ts[i].Ipv4, " ")
		ts[i].Ipv6 = strings.TrimRight(ts[i].Ipv6, " ")
		ts[i].Path = strings.TrimRight(ts[i].Path, " ")
	}
	return ts, nil
}
