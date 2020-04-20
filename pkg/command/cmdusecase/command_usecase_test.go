package cmdusecase_test

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/mojodojo101/c2server/pkg/command/commanddb"

	"github.com/mojodojo101/c2server/pkg/command/cmdusecase"
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
	bu := cmdusecase.NewCommandUsecase(br, time.Second*10)
	err = bu.CreateTable(ctx)
	assert.NoError(t, err)

}

func TestGetNextCommand(t *testing.T) {

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	cr := commanddb.NewSQLRepo(db)
	ctx := context.Background()

	cu := cmdusecase.NewCommandUsecase(cr, time.Second*2)

	TId := int64(1)

	cmd, err := cu.GetNextCommand(ctx, TId)
	fmt.Printf("%v\n", cmd)
	assert.NoError(t, err)

}

func TestListCommandsByTargetID(t *testing.T) {

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	cr := commanddb.NewSQLRepo(db)
	ctx := context.Background()

	cu := cmdusecase.NewCommandUsecase(cr, time.Second*2)

	TId := int64(1)
	amount := int64(20)

	cmd, err := cu.ListCommandsByTargetID(ctx, TId, amount)
	fmt.Printf("%v\n", cmd)
	assert.NoError(t, err)

}
func TestUpdate(t *testing.T) {

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	cr := commanddb.NewSQLRepo(db)
	ctx := context.Background()

	cu := cmdusecase.NewCommandUsecase(cr, time.Second*2)
	c := models.Command{}
	c.Cmd = "start calc.exe"
	c.Id = 1
	c.TId = 1
	c.Executed = true
	c.ExecutedAt = time.Time{}
	c.CreatedAt = time.Now()

	err = cu.Update(ctx, &c)
	assert.NoError(t, err)

}

func TestStore(t *testing.T) {

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	cr := commanddb.NewSQLRepo(db)
	ctx := context.Background()

	cu := cmdusecase.NewCommandUsecase(cr, time.Second*2)
	c := models.Command{}
	c.Cmd = "where calc.exe"
	c.Id = 1
	c.TId = 1
	c.Executed = false
	c.ExecutedAt = time.Time{}
	c.CreatedAt = time.Now()

	err = cu.Store(ctx, &c)
	assert.NoError(t, err)

}
func TestDelete(t *testing.T) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	cr := commanddb.NewSQLRepo(db)
	ctx := context.Background()

	cu := cmdusecase.NewCommandUsecase(cr, time.Second*2)
	c := models.Command{}
	c.Cmd = "start calc.exe"
	c.Id = 1
	err = cu.Delete(ctx, &c)
	assert.NoError(t, err)

}
