package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"github.com/mojodojo101/c2server/pkg/activebeacon/activebeacondb"
	"github.com/mojodojo101/c2server/pkg/beacon/beacondb"
	"github.com/mojodojo101/c2server/pkg/client/clientdb"
	"github.com/mojodojo101/c2server/pkg/command/commanddb"
	"github.com/mojodojo101/c2server/pkg/models"
	"github.com/mojodojo101/c2server/pkg/target/targetdb"
)

func main() {
	for i := 0; i < 5; i++ {
		fmt.Println("")
	}
	fmt.Println("starting test")
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
	ctx := context.Background()
	repoClient := clientdb.NewSQLClientRepo(db)
	err = repoClient.CreateClientTable(ctx)
	if err != nil {
		panic(err)
	}
	clientIp := "127.0.0.1"
	clientName := "mojo"
	clientPassword := "mojodojo101+"
	clientToken := "myawesometoken"
	clientCSRFToken := "myawesometokencsrf"
	clientUpdated_at := time.Now()
	clientCreated_at := time.Now()
	client, err := repoClient.CreateNewClient(ctx,
		clientIp,
		clientName,
		clientPassword,
		clientToken,
		clientCSRFToken,
		clientUpdated_at,
		clientCreated_at)

	fmt.Printf("\n%#v\n", client)

	repoBeacon := beacondb.NewSQLBeaconRepo(db)
	err = repoBeacon.CreateTable(ctx)
	if err != nil {
		panic(err)
	}
	b := models.Beacon{}
	b.Path = "/root/myawesome/beacons/beacon1"
	b.Os = "ubuntu 14.1"
	b.Arch = "x86"
	b.Lang = "PE"

	err = repoBeacon.CreateNewBeacon(ctx, &b)
	fmt.Printf("\n%#v\n", b)

	repoCommand := commanddb.NewSQLCommandRepo(db)
	var targetID int64 = 1
	err = repoCommand.CreateCommandTable(ctx, targetID)
	if err != nil {
		panic(err)
	}

	commandCmd := "start calc.exe"
	commandExecutedAt := time.Now()
	command, err := repoCommand.CreateNewCommand(ctx,
		targetID,
		commandCmd,
		commandExecutedAt)
	fmt.Printf("\n%#v\n", command)

	repoTarget := targetdb.NewSQLTargetRepo(db)
	err = repoTarget.CreateTargetTable(ctx)
	if err != nil {
		panic(err)
	}
	targetIpv4 := "192.168.138.145"
	targetIpv6 := "DEAD::::BEEF"
	targetHostName := "mojo"
	targetPath := "~/c2/out"
	targetCreatedAt := time.Now()
	targetUpadtedAt := time.Now()
	target, err := repoTarget.CreateNewTarget(ctx,
		targetIpv4,
		targetIpv6,
		targetHostName,
		targetPath,
		targetCreatedAt,
		targetUpadtedAt,
	)

	fmt.Printf("\n%#v\n", target)
	ab := models.ActiveBeacon{}
	ab.BId = 0
	ab.C2m = models.HTTP
	ab.TId = target.Id
	ab.CmdId = target.CmdId
	ab.Ping = 0.0
	ab.CreatedAt = time.Now()
	ab.UpdatedAt = time.Now()

	repoActivebeacon := activebeacondb.NewSQLRepo(db)
	err = repoActivebeacon.CreateTable(ctx)
	if err != nil {
		panic(err)
	}
	err = repoActivebeacon.CreateNewBeacon(ctx, &ab)

	if err != nil {
		panic(err)
	}
	fmt.Printf("\n%#v\n", ab)
	//if err = repoBeacon.DeleteById(ctx, beacon.Id); err != nil {
	//	fmt.Printf("\n%v\n", err)
	//}
	//fmt.Println("THIS WORKED FEELSGOODMAN")
}
