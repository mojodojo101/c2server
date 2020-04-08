package activebeacondb

//probably gonna repimplement all of this stuff a bit nicer with comments, but for now this works

import (
	"context"
	"database/sql"

	"github.com/sirupsen/logrus"

	"github.com/mojodojo101/c2server/pkg/activebeacon"
	"github.com/mojodojo101/c2server/pkg/models"
)

type sqlActiveBeaconRepo struct {
	DB *sql.DB
}

//BEGIN struct exported methods

func NewSQLRepo(db *sql.DB) activebeacon.Repository {
	return &sqlActiveBeaconRepo{
		DB: db,
	}
}

func (r *sqlActiveBeaconRepo) CreateTable(ctx context.Context) error {
	query := `CREATE TABLE IF NOT EXISTS active_beacon (
		id BIGSERIAL PRIMARY KEY,
		b_id BIGINT REFERENCES beacon(id),
		p_id BIGINT,
		t_id BIGINT REFERENCES target(id),
		cmd_id BIGINT,
		token CHAR(128),
		cmd CHAR(4096),
		ping FLOAT(53),
		c2m INT,
		pm INT,
		missed_pings INT,
		created_at TIMESTAMP,
		updated_at TIMESTAMP)`

	_, err := r.execQuery(ctx, query)

	return err
}

func (r *sqlActiveBeaconRepo) DropTable(ctx context.Context) error {
	query := `DROP TABLE active_beacon`
	_, err := r.execQuery(ctx, query)

	return err
}

func (r *sqlActiveBeaconRepo) GetByID(ctx context.Context, id int64) (*models.ActiveBeacon, error) {
	query := `SELECT * FROM active_beacon WHERE id=$1`
	return r.getOneItem(ctx, query, id)
}

func (r *sqlActiveBeaconRepo) GetByParentID(ctx context.Context, pId int64) (*models.ActiveBeacon, error) {
	//change this later mojo!!!!!!!
	query := `SELECT * FROM active_beacon WHERE p_id=$1`
	return r.getOneItem(ctx, query, pId)
}

func (r *sqlActiveBeaconRepo) DeleteByID(ctx context.Context, id int64) error {
	//change this later mojo!!!!!!!
	query := `SELECT * FROM active_beacon WHERE id=$1`

	_, err := r.execQuery(ctx, query, id)

	return err

}

func (r *sqlActiveBeaconRepo) Update(ctx context.Context, b *models.ActiveBeacon) error {

	query := `Update active_beacon $1,$2,$3,$4,$5,$6,$7,$7,$8,$9,$10,$11 where id = $12`

	_, err := r.execQuery(ctx, query,
		&b.BId,
		&b.PId,
		&b.TId,
		&b.CmdId,
		&b.Token,
		&b.Cmd,
		&b.Ping,
		&b.C2m,
		&b.Pm,
		&b.MissedPings,
		&b.CreatedAt,
		&b.UpdatedAt,
		&b.Id)

	return err

}
func (r *sqlActiveBeaconRepo) CreateNewBeacon(ctx context.Context, b *models.ActiveBeacon) error {

	query := `INSERT INTO active_beacon VALUES ( nextval('active_beacon_id_seq') ,$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) returning id`

	key, err := r.addItem(
		ctx,
		query,
		&b.BId,
		&b.PId,
		&b.TId,
		&b.CmdId,
		&b.Token,
		&b.Cmd,
		&b.Ping,
		&b.C2m,
		&b.Pm,
		&b.MissedPings,
		&b.CreatedAt,
		&b.UpdatedAt)

	if err != nil {
		logrus.Error(err)
		return err
	}

	b.Id = key

	return nil

}

//END struct exported methods

func (r *sqlActiveBeaconRepo) getOneItem(ctx context.Context, query string, args ...interface{}) (*models.ActiveBeacon, error) {
	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	row := stmt.QueryRowContext(ctx, args...)

	b := &models.ActiveBeacon{}

	err = row.Scan(
		&b.Id,
		&b.BId,
		&b.PId,
		&b.TId,
		&b.CmdId,
		&b.Token,
		&b.Cmd,
		&b.Ping,
		&b.C2m,
		&b.Pm,
		&b.MissedPings,
		&b.CreatedAt,
		&b.UpdatedAt)

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return b, nil
}

func (r *sqlActiveBeaconRepo) addItem(ctx context.Context, query string, args ...interface{}) (int64, error) {
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
func (r *sqlActiveBeaconRepo) execQuery(ctx context.Context, query string, args ...interface{}) (int64, error) {

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
