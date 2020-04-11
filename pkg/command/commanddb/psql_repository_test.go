package commanddb_test

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/mojodojo101/c2server/pkg/command/commanddb"
	"github.com/mojodojo101/c2server/pkg/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var connStr = "host=localhost user=c2admin password=mojodojo101+ dbname=c2db port=5432 sslmode=require"

func TestCreateTable(t *testing.T) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	cr := commanddb.NewSQLRepo(db)
	err = cr.CreateTable(ctx)
	assert.NoError(t, err)
	return
}

func TestDropTable(t *testing.T) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	cr := commanddb.NewSQLRepo(db)
	err = cr.CreateTable(ctx)
	if err != nil {
		panic(err)
	}

	cr.DropTable(ctx)
	assert.NoError(t, err)
	return
}
func TestCreateNewCommand(t *testing.T) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	for i := 0; i < 10; i++ {
		c := models.Command{}
		c.Id = 1
		c.TId = 1
		c.Cmd = "start calc.exe"
		c.Executed = false
		c.CreatedAt = time.Now()
		c.ExecutedAt = time.Time{}
		ctx := context.Background()
		cr := commanddb.NewSQLRepo(db)
		err = cr.CreateTable(ctx)
		if err != nil {
			panic(err)
		}

		err = cr.CreateNewCommand(ctx, &c)
		assert.NoError(t, err)
	}
	return

}
func TestUpdate(t *testing.T) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	c := models.Command{}
	c.Id = 1
	c.TId = 1
	c.Cmd = "start calc.exe"
	c.Executed = false
	c.Executing = true
	c.CreatedAt = time.Now()
	c.ExecutedAt = time.Time{}
	ctx := context.Background()
	cr := commanddb.NewSQLRepo(db)
	err = cr.CreateTable(ctx)
	if err != nil {
		panic(err)
	}

	err = cr.Update(ctx, &c)
	assert.NoError(t, err)
	return

}

func TestGetByID(t *testing.T) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	cr := commanddb.NewSQLRepo(db)
	err = cr.CreateTable(ctx)
	if err != nil {
		panic(err)
	}
	id := int64(1)
	c, err := cr.GetByID(ctx, id)
	assert.NoError(t, err)
	assert.NotEmpty(t, c)
	return
}
func TestGetNextCommandToExecute(t *testing.T) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	cr := commanddb.NewSQLRepo(db)
	err = cr.CreateTable(ctx)
	if err != nil {
		panic(err)
	}
	targetId := int64(1)
	cs, err := cr.GetNextCommand(ctx, targetId)
	assert.NotEmpty(t, cs)
	assert.NoError(t, err)
	return
}

func TestGetByTargetID(t *testing.T) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	cr := commanddb.NewSQLRepo(db)
	err = cr.CreateTable(ctx)
	if err != nil {
		panic(err)
	}
	targetId := int64(1)
	amount := int64(5)
	cs, err := cr.GetByTargetID(ctx, amount, targetId)
	for _, c := range *cs {
		assert.NotEmpty(t, c)
	}
	assert.NoError(t, err)
	return
}
func DeleteByID(t *testing.T) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	cr := commanddb.NewSQLRepo(db)
	err = cr.CreateTable(ctx)
	if err != nil {
		panic(err)
	}
	commandID := int64(1)
	err = cr.DeleteByID(ctx, commandID)

	assert.NoError(t, err)
	return
}
