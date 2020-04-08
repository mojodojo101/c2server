package activebeacondb_test

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/mojodojo101/c2server/pkg/activebeacon/activebeacondb"
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
	br := activebeacondb.NewSQLRepo(db)
	err = br.CreateTable(ctx)
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
	br := activebeacondb.NewSQLRepo(db)
	err = br.CreateTable(ctx)
	if err != nil {
		panic(err)
	}
	br.DropTable(ctx)
	assert.NoError(t, err)
	return
}
func TestCreateNewActiveBeacon(t *testing.T) {
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
	br := activebeacondb.NewSQLRepo(db)
	err = br.CreateTable(ctx)
	if err != nil {
		panic(err)
	}
	ab := models.ActiveBeacon{}
	ab.BId = 0
	ab.C2m = models.HTTP
	ab.TId = 1
	ab.CmdId = 1
	ab.Ping = 0.0
	ab.CreatedAt = time.Now()
	ab.UpdatedAt = time.Now()

	err = br.CreateNewBeacon(ctx, &ab)
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
	br := activebeacondb.NewSQLRepo(db)
	err = br.CreateTable(ctx)
	if err != nil {
		panic(err)
	}
	id := int64(1)
	b, err := br.GetByID(ctx, id)
	assert.NoError(t, err)
	assert.NotEmpty(t, b)
	return
}
func TestDeleteByID(t *testing.T) {
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
	br := activebeacondb.NewSQLRepo(db)
	err = br.CreateTable(ctx)
	if err != nil {
		panic(err)
	}
	id := int64(1)
	err = br.DeleteByID(ctx, id)
	assert.NoError(t, err)
}
