package cmdusecase_test

import (
	"context"
	"database/sql"
	"github.com/mojodojo101/c2server/pkg/command/commanddb"

	"github.com/mojodojo101/c2server/pkg/command/usecase"
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
	br := commanddb.NewSQLRepo(db)
	ctx := context.Background()
	bu := usecase.NewCommandUsecase(br, time.Second*10)
	err = bu.CreateTable(ctx)
	assert.NoError(t, err)

}

func TestStore(t *testing.T) {

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	cr := commanddb.NewSQLRepo(db)
	ctx := context.Background()

	cu := usecase.NewCommandUsecase(cr, time.Second*2)
	c := models.Command{}
	c.Cmd = "start calc.exe"
	c.Id = 1
	c.TId = 1
	c.Executed = true
	c.ExecutedAt = time.Time{}
	c.CreatedAt = time.Now()

	_, err = cu.Store(ctx, &c)
	assert.NoError(t, err)

}
func TestDelete(t *testing.T) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	cr := commanddb.NewSQLRepo(db)
	ctx := context.Background()

	cu := usecase.NewCommandUsecase(cr, time.Second*2)
	c := models.Command{}
	c.Cmd = "start calc.exe"
	c.Id = 1
	err = cu.Delete(ctx, &c)
	assert.NoError(t, err)

}
