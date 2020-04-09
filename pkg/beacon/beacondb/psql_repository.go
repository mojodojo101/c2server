package beacondb

import (
	"context"
	"database/sql"

	"github.com/sirupsen/logrus"

	"github.com/mojodojo101/c2server/pkg/beacon"
	"github.com/mojodojo101/c2server/pkg/models"
)

type sqlBeaconRepo struct {
	DB *sql.DB
}

//BEGIN struct exported methods

func NewSQLRepo(db *sql.DB) beacon.Repository {
	return &sqlBeaconRepo{
		DB: db,
	}
}

func (r *sqlBeaconRepo) CreateTable(ctx context.Context) error {
	query := `CREATE TABLE IF NOT EXISTS beacon (
		id BIGSERIAL PRIMARY KEY,
		path CHAR(255),
		os CHAR(100),
		arch CHAR(20),
		lang CHAR(100))`

	_, err := r.execQuery(ctx, query)

	return err
}

func (r *sqlBeaconRepo) DropTable(ctx context.Context) error {
	query := `DROP TABLE beacon`

	_, err := r.execQuery(ctx, query)

	return err
}
func (r *sqlBeaconRepo) GetByID(ctx context.Context, id int64) (*models.Beacon, error) {
	query := `SELECT * FROM beacon WHERE id=$1`
	return r.getOneItem(ctx, query, id)
}

func (r *sqlBeaconRepo) DeleteByID(ctx context.Context, id int64) error {
	query := `DELETE FROM beacon WHERE id=$1`
	_, err := r.execQuery(ctx, query, id)
	return err
}

func (r *sqlBeaconRepo) CreateNewBeacon(ctx context.Context, b *models.Beacon) (int64, error) {
	query := `INSERT INTO beacon VALUES ( nextval('beacon_id_seq') ,$1,$2,$3,$4) returning id`

	key, err := r.addItem(
		ctx,
		query,
		&b.Path,
		&b.Os,
		&b.Arch,
		&b.Lang,
	)

	if err != nil {
		logrus.Error(err)
		return int64(0), err
	}

	return key, nil

}

//END struct exported methods

func (r *sqlBeaconRepo) getOneItem(ctx context.Context, query string, args ...interface{}) (*models.Beacon, error) {
	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	row := stmt.QueryRowContext(ctx, args...)

	b := &models.Beacon{}

	err = row.Scan(
		&b.Id,
		&b.Path,
		&b.Os,
		&b.Arch,
		&b.Lang,
	)

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return b, nil
}

func (r *sqlBeaconRepo) addItem(ctx context.Context, query string, args ...interface{}) (int64, error) {
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
func (r *sqlBeaconRepo) execQuery(ctx context.Context, query string, args ...interface{}) (int64, error) {

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
