package abusecase_test

import (
	"context"
	"database/sql"
	"github.com/mojodojo101/c2server/pkg/activebeacon/activebeacondb"

	"github.com/mojodojo101/c2server/pkg/activebeacon/abusecase"
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
	ar := activebeacondb.NewSQLRepo(db)
	ctx := context.Background()
	au := usecase.NewActiveBeaconUsecase(ar, time.Second*10)
	err = au.CreateTable(ctx)
	assert.NoError(t, err)

}

func TestRegister(t *testing.T) {

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	ar := activebeacondb.NewSQLRepo(db)
	ctx := context.Background()
	au := usecase.NewActiveBeaconUsecase(ar, time.Second*10)

	a := models.ActiveBeacon{}
	a.Id = 1
	a.PId = 0
	a.BId = 1
	a.Cmd = "start calc.exe"
	a.CmdId = 1
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()
	a.TId = 1
	a.MissedPings = 0
	a.Pm = models.HTTP
	a.C2m = models.HTTP
	a.Token = "23910809213"
	a.Ping = float64(10.0)

	_, err = au.Register(ctx, &a)
	assert.NoError(t, err)

}
func TestUpdate(t *testing.T) {

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	ar := activebeacondb.NewSQLRepo(db)
	ctx := context.Background()
	au := usecase.NewActiveBeaconUsecase(ar, time.Second*10)

	a := models.ActiveBeacon{}
	a.Id = 1
	a.PId = 0
	a.BId = 1
	a.Cmd = "start calc.exe"
	a.CmdId = 1
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()
	a.TId = 1
	a.MissedPings = 0
	a.Pm = models.HTTP
	a.C2m = models.HTTP
	a.Token = "23910809213"
	a.Ping = float64(10.0)

	err = au.Update(ctx, &a)
	assert.NoError(t, err)

}
