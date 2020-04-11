package commanddb

import (
	"context"
	"database/sql"

	"github.com/sirupsen/logrus"

	"github.com/mojodojo101/c2server/pkg/command"
	"github.com/mojodojo101/c2server/pkg/models"
)

type sqlCommandRepo struct {
	DB *sql.DB
}

//BEGIN struct exported methods

//there is probably a better way to generate these table names, if u (whoever reads this has any suggestions please let me know (pms are open) https://twitter.com/Mojodojo_101
//i might just refactor all this, not sure yet, this is a learning experience for me, could also just have 1 table for all the commands

func NewSQLRepo(db *sql.DB) command.Repository {
	return &sqlCommandRepo{
		DB: db,
	}
}

func (r *sqlCommandRepo) CreateTable(ctx context.Context) error {

	query := `CREATE TABLE IF NOT EXISTS command (
		id BIGSERIAL PRIMARY KEY,
		t_id BIGINT,
		cmd char(4096),
		executed BOOL,
		executing BOOL,
		created_at TIMESTAMP,
		executed_at TIMESTAMP)`
	_, err := r.execQuery(ctx, query)

	return err
}

func (r *sqlCommandRepo) DropTable(ctx context.Context) error {
	query := `DROP TABLE command`
	_, err := r.execQuery(ctx, query)

	return err
}

func (r *sqlCommandRepo) GetByID(ctx context.Context, id int64) (*models.Command, error) {
	query := `SELECT * FROM command WHERE id=$1`
	return r.getOneItem(ctx, query, id)
}

func (r *sqlCommandRepo) GetNextCommand(ctx context.Context, targetId int64) (*models.Command, error) {
	query := `SELECT * FROM command WHERE t_id=$1 and executing=FALSE`
	return r.getOneItem(ctx, query, targetId)
}
func (r *sqlCommandRepo) GetByTargetID(ctx context.Context, amount, targetId int64) (*[]models.Command, error) {
	query := `SELECT * FROM command WHERE t_id=$1`
	return r.getItemsByValue(ctx, query, amount, targetId)

}

func (r *sqlCommandRepo) DeleteByID(ctx context.Context, id int64) error {
	query := `DELETE from command WHERE id=$1`
	_, err := r.execQuery(ctx, query, id)
	return err
}

//i should change this stuff down the line
func (r *sqlCommandRepo) Update(ctx context.Context, c *models.Command) error {
	query := `Update command SET executing=$2 ,executed = $3, executed_at = $4 WHERE id=$1`
	_, err := r.execQuery(ctx, query, c.Id, c.Executing, c.Executed, c.ExecutedAt)
	return err
}

func (r *sqlCommandRepo) CreateNewCommand(ctx context.Context, c *models.Command) error {
	query := `INSERT INTO command VALUES ( nextval('command_id_seq') ,$1,$2,$3,$4,$5,$6) returning id`

	key, err := r.addItem(
		ctx,
		query,
		&c.TId,
		&c.Cmd,
		&c.Executed,
		&c.Executing,
		&c.CreatedAt,
		&c.ExecutedAt,
	)

	if err != nil {
		logrus.Error(err)
	}
	c.Id = key
	return err

}

//END struct exported methods

func (r *sqlCommandRepo) getOneItem(ctx context.Context, query string, args ...interface{}) (*models.Command, error) {
	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	row := stmt.QueryRowContext(ctx, args...)

	c := &models.Command{}

	err = row.Scan(
		&c.Id,
		&c.TId,
		&c.Cmd,
		&c.Executed,
		&c.Executing,
		&c.CreatedAt,
		&c.ExecutedAt,
	)

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return c, nil
}

func (r *sqlCommandRepo) getItemsByValue(ctx context.Context, query string, amount int64, args ...interface{}) (*[]models.Command, error) {
	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	row := stmt.QueryRowContext(ctx, args...)
	cs := make([]models.Command, amount)
	for i, _ := range cs {
		err = row.Scan(
			&cs[i].Id,
			&cs[i].TId,
			&cs[i].Cmd,
			&cs[i].Executed,
			&cs[i].Executing,
			&cs[i].CreatedAt,
			&cs[i].ExecutedAt,
		)

		if err != nil {
			logrus.Error(err)
			if i != 0 {
				cs = append(cs[:i])
				return &cs, nil
			}
			return nil, err
		}
	}
	return &cs, nil
}

func (r *sqlCommandRepo) addItem(ctx context.Context, query string, args ...interface{}) (int64, error) {
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
func (r *sqlCommandRepo) execQuery(ctx context.Context, query string, args ...interface{}) (int64, error) {

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
