package abusecase_test

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/mojodojo101/c2server/pkg/activebeacon/activebeacondb"

	"github.com/mojodojo101/c2server/pkg/activebeacon/abusecase"
	"github.com/mojodojo101/c2server/pkg/command/cmdusecase"
	"github.com/mojodojo101/c2server/pkg/command/commanddb"
	"github.com/mojodojo101/c2server/pkg/target/targetdb"
	"github.com/mojodojo101/c2server/pkg/target/tusecase"
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
	tr := targetdb.NewSQLRepo(db)
	cmdr := commanddb.NewSQLRepo(db)
	ctx := context.Background()

	cu := cmdusecase.NewCommandUsecase(cmdr, time.Second*2)
	tu := tusecase.NewTargetUsecase(tr, cu, time.Second*2)
	au := abusecase.NewActiveBeaconUsecase(ar, tu, time.Second*2)

	err = au.CreateTable(ctx)
	assert.NoError(t, err)

}

func TestRegister(t *testing.T) {

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	ar := activebeacondb.NewSQLRepo(db)
	tr := targetdb.NewSQLRepo(db)
	cmdr := commanddb.NewSQLRepo(db)
	ctx := context.Background()

	cu := cmdusecase.NewCommandUsecase(cmdr, time.Second*2)
	tu := tusecase.NewTargetUsecase(tr, cu, time.Second*2)
	au := abusecase.NewActiveBeaconUsecase(ar, tu, time.Second*2)

	a := models.ActiveBeacon{}
	a.Id = 321321
	a.PId = 0
	a.BId = 1
	a.CmdId = 1
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()
	a.TId = 1
	a.MissedPings = 0
	a.Pm = models.HTTP
	a.C2m = models.HTTP
	a.Token = "23910809213"
	a.Ping = float64(10.0)

	err = au.Register(ctx, &a)
	fmt.Printf("a.Id = %v\n", a.Id)
	assert.NoError(t, err)

}
func TestSetCmdExecuted(t *testing.T) {

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	ar := activebeacondb.NewSQLRepo(db)
	tr := targetdb.NewSQLRepo(db)
	cmdr := commanddb.NewSQLRepo(db)
	ctx := context.Background()

	cu := cmdusecase.NewCommandUsecase(cmdr, time.Second*2)
	tu := tusecase.NewTargetUsecase(tr, cu, time.Second*2)
	au := abusecase.NewActiveBeaconUsecase(ar, tu, time.Second*2)

	a := models.ActiveBeacon{}
	a.Id = 1
	a.PId = 0
	a.BId = 1
	a.CmdId = 1
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()
	a.TId = 1
	a.MissedPings = 0
	a.Pm = models.HTTP
	a.C2m = models.HTTP
	a.Token = "mytoken"
	a.Ping = float64(10.0)

	err = au.SetCmdExecuted(ctx, &a, []byte("some response from the network"))

	assert.NoError(t, err)

}
func TestUpdate(t *testing.T) {

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	ar := activebeacondb.NewSQLRepo(db)
	tr := targetdb.NewSQLRepo(db)
	cmdr := commanddb.NewSQLRepo(db)
	ctx := context.Background()

	cu := cmdusecase.NewCommandUsecase(cmdr, time.Second*2)
	tu := tusecase.NewTargetUsecase(tr, cu, time.Second*2)
	au := abusecase.NewActiveBeaconUsecase(ar, tu, time.Second*2)

	a := models.ActiveBeacon{}
	a.Id = 1
	a.PId = 0
	a.BId = 1
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
