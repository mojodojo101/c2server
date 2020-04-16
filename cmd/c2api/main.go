package main

import (
	"context"
	"database/sql"
	_ "fmt"
	"net/http"
	"time"

	_ "github.com/lib/pq"

	"github.com/mojodojo101/c2server/pkg/activebeacon/abusecase"
	"github.com/mojodojo101/c2server/pkg/activebeacon/activebeacondb"
	"github.com/mojodojo101/c2server/pkg/activebeacon/delivery/abhttp"

	_ "github.com/mojodojo101/c2server/pkg/beacon/beacondb"
	_ "github.com/mojodojo101/c2server/pkg/beacon/busecase"
	"github.com/mojodojo101/c2server/pkg/client/clientdb"
	"github.com/mojodojo101/c2server/pkg/client/cusecase"
	"github.com/mojodojo101/c2server/pkg/client/delivery/chttp"
	"github.com/mojodojo101/c2server/pkg/command/cmdusecase"
	"github.com/mojodojo101/c2server/pkg/command/commanddb"
	_ "github.com/mojodojo101/c2server/pkg/models"
	"github.com/mojodojo101/c2server/pkg/target/targetdb"
	"github.com/mojodojo101/c2server/pkg/target/tusecase"
)

func main() {
	connStr := "host=localhost user=c2admin password=mojodojo101+ dbname=c2db port=5432 sslmode=require"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	timeout := time.Second * 6

	//init client repo
	ctx := context.Background()
	//init beacon repo
	//br := beacondb.NewSQLRepo(db)
	//err = br.CreateTable(ctx)
	//if err != nil {
	//	panic(err)
	//}
	//init beacon usecase
	///bu := busecase.NewBeaconUsecase(br, timeout)
	///if err != nil {
	///	panic(err)
	///}

	//init command repo
	cmdr := commanddb.NewSQLRepo(db)
	err = cmdr.CreateTable(ctx)
	if err != nil {
		panic(err)
	}
	//init command usecase
	cmdu := cmdusecase.NewCommandUsecase(cmdr, timeout)
	if err != nil {
		panic(err)
	}

	//init target repo
	tr := targetdb.NewSQLRepo(db)
	err = tr.CreateTable(ctx)
	if err != nil {
		panic(err)
	}
	//init target usecase
	tu := tusecase.NewTargetUsecase(tr, cmdu, timeout)
	if err != nil {
		panic(err)
	}

	//init activebeacon repo
	ar := activebeacondb.NewSQLRepo(db)
	err = ar.CreateTable(ctx)
	if err != nil {
		panic(err)
	}
	//init activebeacon usecase
	au := abusecase.NewActiveBeaconUsecase(ar, tu, timeout)
	if err != nil {
		panic(err)
	}

	//init client repo
	cr := clientdb.NewSQLRepo(db)
	err = cr.CreateTable(ctx)
	if err != nil {
		panic(err)
	}
	//init client usecase
	cu := cusecase.NewClientUsecase(cr, tu, timeout)
	if err != nil {
		panic(err)
	}

	ch := chttp.NewHandler(cu)
	abh := abhttp.NewHandler(au)
	signalCh := make(chan bool, 1)
	go func() {
		for {
			select {
			case <-signalCh:
				return
			default:
				http.ListenAndServe(":80", &abh)
			}
		}
	}()

	http.ListenAndServe(":443", &ch)
	signalCh <- true
}
