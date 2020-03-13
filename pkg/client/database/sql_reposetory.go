package database

import(
	"context"
	"time"
	"database/sql"

	"github.com/sirupsen/logrus"

	"github.com/mojodojo101/c2server/pkg/client"
	"github.com/mojodojo101/c2server/pkg/models"
)

type sqlClientRepo struct {
	DB *sql.DB
}


func NewSQLClientRepo(db *sql.DB) client.Repository {
	return &sqlClientRepo{
		DB: db,
	}
}

//BEGIN struct exported methods

func (r *sqlClientRepo) CreateClientTable(ctx context.Context)(error){
	query := `CREATE TABLE IF NOT EXISTS client (
		id BIGSERIAL PRIMARY KEY,
		ip CHAR(16),
		name CHAR(50),
		password CHAR(128),
		token CHAR(128),
		created_at TIMESTAMP,
		updated_at TIMESTAMP)`

	_,err:=r.execQuery(ctx,query)

	return err
}


func (r *sqlClientRepo) GetByID(ctx context.Context, id int64) (*models.Client, error) {
	query := `SELECT id,ip,name,password,token, created_at, updated_at FROM author WHERE id=$1`
	return r.getOneItem(ctx, query, id)
}


func  (r *sqlClientRepo) CreateNewClient(ctx context.Context,ip,name,password,token string,createdAt,updatedAt time.Time)(*models.Client,error){
	c := models.Client{
		Id : 0,
		Ip :ip,
		Name :name,
		Password :password,
		Token : token,
		CreatedAt	: createdAt,
		UpdatedAt	: updatedAt,
	}

	query:=`INSERT INTO client VALUES ( nextval('client_id_seq') ,$1,$2,$3,$4,$5,$6) returning id`

	key,err:=r.addItem(
		ctx,
		query,
		&c.Ip,
		&c.Name,
		&c.Password,
		&c.Token,
		&c.CreatedAt,
		&c.UpdatedAt)
	

    if err!=nil{
		logrus.Error(err)
		return nil,err
	}

	c.Id=key

	return &c,nil


}




//END struct exported methods

func (r *sqlClientRepo) getOneItem(ctx context.Context,query string,args ...interface{})(*models.Client ,error){
	stmt,err:= r.DB.PrepareContext(ctx,query)
	if err != nil {
		logrus.Error(err)
		return nil,err
	}

	row:=stmt.QueryRowContext(ctx,args...)

	c:=&models.Client{}

	err = row.Scan(
		&c.Id,
		&c.Ip,
		&c.Name,
		&c.Password,
		&c.Token,
		&c.CreatedAt,
		&c.UpdatedAt,
	)

	if err !=nil {
		logrus.Error(err)
		return nil,err
	}

	return c,nil
}

func (r *sqlClientRepo) addItem(ctx context.Context,query string,args ...interface{})(int64,error){
	tx, err := r.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		logrus.Error(err)
		return 0,err
	}
	var id int64

	err =tx.QueryRowContext(ctx,query,args...).Scan(&id)

	if err!= nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			logrus.Error(err)
		}
		logrus.Error(err)
	}

	if err := tx.Commit(); err != nil {
		logrus.Error(err)
	}

	return id,nil

}

// executes command and returns rows affected and error
func (r *sqlClientRepo) execQuery(ctx context.Context,query string,args ...interface{})(int64,error){

	tx, err := r.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		logrus.Error(err)
		return 0,err
	}


	result,err :=tx.ExecContext(ctx,query,args...)
	if err!= nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			logrus.Error(err)
		}
		logrus.Error(err)
	}

	if err := tx.Commit(); err != nil {
		logrus.Error(err)
	}
	rows,err:=result.RowsAffected()

	return rows,nil

}