package clientdb

//probably gonna repimplement all of this stuff a bit nicer with comments, but for now this works

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/mojodojo101/c2server/pkg/client"
	"github.com/mojodojo101/c2server/pkg/models"
)

type sqlClientRepo struct {
	DB *sql.DB
}

//BEGIN struct exported methods

func NewSQLRepo(db *sql.DB) client.Repository {
	return &sqlClientRepo{
		DB: db,
	}
}

func (r *sqlClientRepo) CreateTable(ctx context.Context) error {
	query := `CREATE TABLE IF NOT EXISTS client (
		id BIGSERIAL PRIMARY KEY,
		ip CHAR(16),
		name CHAR(50),
		password CHAR(128),
		token CHAR(128),
		csrf_token CHAR(128),
		created_at TIMESTAMP,
		updated_at TIMESTAMP,
		UNIQUE(token,name))
		`

	_, err := r.execQuery(ctx, query)

	return err
}

func (r *sqlClientRepo) DropTable(ctx context.Context) error {
	query := `DROP TABLE client`
	_, err := r.execQuery(ctx, query)

	return err
}

func (r *sqlClientRepo) GetByID(ctx context.Context, id int64) (*models.Client, error) {
	query := `SELECT * FROM client WHERE id=$1`
	return r.getOneItem(ctx, query, id)
}
func (r *sqlClientRepo) GetByNameAndPassword(ctx context.Context, name, password string) (*models.Client, error) {
	query := `SELECT * FROM client WHERE name=$1 and password=$2`
	return r.getOneItem(ctx, query, name, password)
}
func (r *sqlClientRepo) DeleteByID(ctx context.Context, id int64) error {
	query := `DELETE FROM client WHERE id=$1`
	_, err := r.execQuery(ctx, query, id)
	return err
}
func (r *sqlClientRepo) UpdateCSRFToken(ctx context.Context, id int64, csrfToken string, updatedAt time.Time) error {
	query := `UPDATE client SET csrf_token= $2 ,updated_at = $3 WHERE id=$1`
	_, err := r.execQuery(ctx, query, id, csrfToken, updatedAt)
	return err
}
func (r *sqlClientRepo) Update(ctx context.Context, c *models.Client) error {
	query := `UPDATE client SET ip=$2,name=$3,password=$4,token=$5,csrf_token=$6,updated_at=$7 WHERE id=$1`
	_, err := r.execQuery(ctx, query, c.Id, c.Ip, c.Name, c.Password, c.Token, c.CSRFToken, c.UpdatedAt)
	return err
}
func (r *sqlClientRepo) CreateNewClient(ctx context.Context, c *models.Client) error {
	query := `INSERT INTO client VALUES ( nextval('client_id_seq') ,$1,$2,$3,$4,$5,$6,$7) returning id`

	key, err := r.addItem(
		ctx,
		query,
		&c.Ip,
		&c.Name,
		&c.Password,
		&c.Token,
		&c.CSRFToken,
		&c.CreatedAt,
		&c.UpdatedAt)

	if err != nil {
		logrus.Error(err)
	}
	c.Name = strings.TrimSpace(c.Name)
	c.Password = strings.TrimSpace(c.Password)
	c.Token = strings.TrimSpace(c.Token)
	c.CSRFToken = strings.TrimSpace(c.CSRFToken)
	c.Id = key

	return err

}

//END struct exported methods

func (r *sqlClientRepo) getOneItem(ctx context.Context, query string, args ...interface{}) (*models.Client, error) {
	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	row := stmt.QueryRowContext(ctx, args...)

	c := &models.Client{}

	err = row.Scan(
		&c.Id,
		&c.Ip,
		&c.Name,
		&c.Password,
		&c.Token,
		&c.CSRFToken,
		&c.CreatedAt,
		&c.UpdatedAt,
	)

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return c, nil
}

func (r *sqlClientRepo) addItem(ctx context.Context, query string, args ...interface{}) (int64, error) {
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
func (r *sqlClientRepo) execQuery(ctx context.Context, query string, args ...interface{}) (int64, error) {

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
