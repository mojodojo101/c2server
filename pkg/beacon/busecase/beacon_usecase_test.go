package busecase_test

import (
	"context"
	"database/sql"
	"github.com/mojodojo101/c2server/pkg/beacon/beacondb"

	"github.com/mojodojo101/c2server/pkg/beacon/busecase"
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
	br := beacondb.NewSQLRepo(db)
	ctx := context.Background()
	bu := usecase.NewBeaconUsecase(br, time.Second*10)
	err = bu.CreateTable(ctx)
	assert.NoError(t, err)

}

func TestStore(t *testing.T) {

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	br := beacondb.NewSQLRepo(db)
	ctx := context.Background()

	bu := usecase.NewBeaconUsecase(br, time.Second*2)
	b := models.Beacon{}
	b.Arch = "x86"
	b.Id = 1
	b.Lang = "bin"
	b.Path = "/root/go/src/github.com/mojodojo101/c2server/internal_resources/beacons/callHome32.dll"
	_, err = bu.Store(ctx, &b)
	assert.Error(t, err)

}
func TestRetrieveBeaconBuffer(t *testing.T) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	br := beacondb.NewSQLRepo(db)
	ctx := context.Background()

	bu := usecase.NewBeaconUsecase(br, time.Second*2)
	b := models.Beacon{}
	b.Arch = "x86"
	b.Id = 1
	b.Lang = "asm"
	b.Path = "/root/go/src/github.com/mojodojo101/c2server/internal_resources/beacons/callHome32.dll"
	buffer, err := bu.RetrieveBeaconBuffer(ctx, &b)
	assert.NoError(t, err)
	assert.NotEmpty(t, buffer)

}
