package cusecase_test

import (
	"context"
	"database/sql"
	"github.com/mojodojo101/c2server/pkg/client/clientdb"

	"github.com/mojodojo101/c2server/pkg/client/usecase"
	"github.com/stretchr/testify/assert"

	_ "github.com/lib/pq"
	"github.com/mojodojo101/c2server/pkg/models"
	"testing"
	"time"
)

var connStr = "host=localhost user=c2admin password=mojodojo101+ dbname=c2db port=5432 sslmode=require"

func TestCreateTable(t *testing.T) {

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	cr := clientdb.NewSQLRepo(db)
	ctx := context.Background()
	cu := usecase.NewClientUsecase(cr, time.Second*10)
	err = cu.CreateTable(ctx)
	assert.NoError(t, err)

}

func TestStore(t *testing.T) {

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	cr := clientdb.NewSQLRepo(db)
	ctx := context.Background()
	cu := usecase.NewClientUsecase(cr, time.Second*10)
	c := models.Client{}
	c.Id = 1
	c.Password = "mypass"
	c.Token = "ndwsioanwoqinoid"
	c.CSRFToken = "32132132121"
	c.Ip = "192.168.122.102"
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	_, err = cu.Store(ctx, &c)
	assert.NoError(t, err)

}
func TestUpdate(t *testing.T) {

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	cr := clientdb.NewSQLRepo(db)
	ctx := context.Background()
	cu := usecase.NewClientUsecase(cr, time.Second*10)
	c := models.Client{}
	c.Id = 1
	c.Password = "mypass"
	c.Token = "this should be different"
	c.CSRFToken = "32132132121"
	c.Ip = "192.168.122.102"
	c.UpdatedAt = time.Now()
	c.CreatedAt = time.Time{}

	err = cu.Update(ctx, &c)
	assert.NoError(t, err)

}
