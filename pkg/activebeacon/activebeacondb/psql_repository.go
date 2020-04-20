package activebeacondb

//probably gonna repimplement all of this stuff a bit nicer with comments, but for now this works

import (
	"context"
	"database/sql"
	"strings"

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
		cmd CHAR(4096),
		token CHAR(128),
		ping BIGINT,
		c2m INT,
		pm INT,
		missed_pings BIGINT,
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
	query := `SELECT * FROM active_beacon WHERE p_id=$1`
	return r.getOneItem(ctx, query, pId)
}
func (r *sqlActiveBeaconRepo) GetAllActiveBeacons(ctx context.Context, amount int64) ([]models.ActiveBeacon, error) {
	query := `SELECT * FROM active_beacon ORDER BY id DESC`
	return r.getAllItems(ctx, query, amount)
}
func (r *sqlActiveBeaconRepo) DeleteByID(ctx context.Context, id int64) error {
	query := `DELETE FROM active_beacon WHERE id=$1`

	_, err := r.execQuery(ctx, query, id)

	return err

}

//simply updates all values i might want to change this later
func (r *sqlActiveBeaconRepo) Update(ctx context.Context, b *models.ActiveBeacon) error {

	query := `Update active_beacon set b_id=$1,p_id=$2,t_id=$3,cmd_id=$4,cmd=$5,token=$6,ping=$7,c2m=$8,pm=$9,missed_pings=$10,created_at=$11,updated_at=$12 where id = $13`

	_, err := r.execQuery(ctx, query,
		&b.BId,
		&b.PId,
		&b.TId,
		&b.CmdId,
		&b.Cmd,
		&b.Token,
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
		&b.Cmd,
		&b.Token,
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

	ab := &models.ActiveBeacon{}

	err = row.Scan(
		&ab.Id,
		&ab.BId,
		&ab.PId,
		&ab.TId,
		&ab.CmdId,
		&ab.Cmd,
		&ab.Token,
		&ab.Ping,
		&ab.C2m,
		&ab.Pm,
		&ab.MissedPings,
		&ab.CreatedAt,
		&ab.UpdatedAt)

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	ab.Cmd = strings.TrimRight(ab.Cmd, " ")
	ab.Token = strings.TrimRight(ab.Token, " ")

	return ab, nil
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

func (r *sqlActiveBeaconRepo) getAllItems(ctx context.Context, query string, amount int64, args ...interface{}) ([]models.ActiveBeacon, error) {
	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer rows.Close()
	abs := make([]models.ActiveBeacon, amount)
	for i, _ := range abs {
		rows.Next()
		err = rows.Scan(
			&abs[i].Id,
			&abs[i].BId,
			&abs[i].PId,
			&abs[i].TId,
			&abs[i].CmdId,
			&abs[i].Cmd,
			&abs[i].Token,
			&abs[i].Ping,
			&abs[i].C2m,
			&abs[i].Pm,
			&abs[i].MissedPings,
			&abs[i].CreatedAt,
			&abs[i].UpdatedAt,
		)

		if err != nil {
			logrus.Error(err)
			if i != 0 {
				abs = append(abs[:i])
				return abs, nil
			}
			return nil, err
		}

		abs[i].Cmd = strings.TrimRight(abs[i].Cmd, " ")
		abs[i].Token = strings.TrimRight(abs[i].Token, " ")
	}
	return abs, nil
}
