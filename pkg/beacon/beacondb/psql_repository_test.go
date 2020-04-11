package beacondb_test

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/mojodojo101/c2server/pkg/beacon/beacondb"
	"github.com/mojodojo101/c2server/pkg/models"
	"github.com/stretchr/testify/assert"
	"testing"
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
	br := beacondb.NewSQLRepo(db)
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
	br := beacondb.NewSQLRepo(db)
	err = br.CreateTable(ctx)
	if err != nil {
		panic(err)
	}
	br.DropTable(ctx)
	assert.NoError(t, err)
	return
}
func TestCreateNewBeacon(t *testing.T) {
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
	br := beacondb.NewSQLRepo(db)
	err = br.CreateTable(ctx)
	if err != nil {
		panic(err)
	}
	b := models.Beacon{}
	b.Path = "/root/myawesome/beacons/beacon1"
	b.Os = "ubuntu 14.1"
	b.Arch = "x86"
	b.Lang = "PE"

	err = br.CreateNewBeacon(ctx, &b)
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
	br := beacondb.NewSQLRepo(db)
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
	br := beacondb.NewSQLRepo(db)
	err = br.CreateTable(ctx)
	if err != nil {
		panic(err)
	}
	id := int64(1)
	err = br.DeleteByID(ctx, id)
	assert.NoError(t, err)
}
