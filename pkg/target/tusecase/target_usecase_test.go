package tusecase_test

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/mojodojo101/c2server/pkg/command/commanddb"
	"github.com/mojodojo101/c2server/pkg/target/targetdb"

	"github.com/mojodojo101/c2server/pkg/command/cmdusecase"
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
	tr := targetdb.NewSQLRepo(db)
	cr := commanddb.NewSQLRepo(db)
	cu := cmdusecase.NewCommandUsecase(cr, time.Second*2)

	ctx := context.Background()
	tu := tusecase.NewTargetUsecase(tr, cu, time.Second*10)
	err = tu.CreateTable(ctx)
	assert.NoError(t, err)

}

func TestStore(t *testing.T) {

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	tr := targetdb.NewSQLRepo(db)
	cr := commanddb.NewSQLRepo(db)
	cu := cmdusecase.NewCommandUsecase(cr, time.Second*2)

	ctx := context.Background()
	tu := tusecase.NewTargetUsecase(tr, cu, time.Second*10)

	ta := models.Target{}
	ta.HostName = "the beast"
	ta.Id = 3
	ta.Ipv4 = "localhost"
	ta.Ipv6 = "DEAD:BEEF::::"
	ta.CreatedAt = time.Now()
	ta.UpdatedAt = time.Now()
	ta.Path = fmt.Sprintf("/root/go/src/github.com/mojodojo101/c2server/internal_resources/targets/%v/", ta.Id)
	err = tu.Store(ctx, &ta)
	assert.NoError(t, err)

}
func TestUpdate(t *testing.T) {

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	tr := targetdb.NewSQLRepo(db)
	cr := commanddb.NewSQLRepo(db)
	cu := cmdusecase.NewCommandUsecase(cr, time.Second*2)

	ctx := context.Background()
	tu := tusecase.NewTargetUsecase(tr, cu, time.Second*10)

	ta := models.Target{}
	ta.HostName = "the beast"
	ta.Id = 1
	ta.Ipv4 = "localhost"
	ta.Ipv6 = "DEAD:BEEF::::"
	ta.CreatedAt = time.Now()
	ta.UpdatedAt = time.Now()
	ta.Path = fmt.Sprintf("/root/go/src/github.com/mojodojo101/c2server/internal_resources/targets/%v/", ta.Id)
	err = tu.Update(ctx, &ta)
	assert.NoError(t, err)

}

func TestGetNextCmd(t *testing.T) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	tr := targetdb.NewSQLRepo(db)
	cr := commanddb.NewSQLRepo(db)
	ctx := context.Background()
	cu := cmdusecase.NewCommandUsecase(cr, time.Second*2)
	tu := tusecase.NewTargetUsecase(tr, cu, time.Second*2)

	ta := models.Target{}
	ta.HostName = "the beast"
	ta.Id = 1
	ta.Ipv4 = "localhost"
	ta.Ipv6 = "DEAD:BEEF::::"
	ta.CreatedAt = time.Now()
	ta.UpdatedAt = time.Now()
	ta.Path = fmt.Sprintf("/root/go/src/github.com/mojodojo101/c2server/internal_resources/targets/%v/", ta.Id)
	cmd, err := tu.GetNextCmd(ctx, &ta)
	assert.NoError(t, err)
	assert.NotEmpty(t, cmd)

}
func TestSetCmdExecuted(t *testing.T) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	tr := targetdb.NewSQLRepo(db)
	cr := commanddb.NewSQLRepo(db)
	ctx := context.Background()
	cu := cmdusecase.NewCommandUsecase(cr, time.Second*2)
	tu := tusecase.NewTargetUsecase(tr, cu, time.Second*2)

	ta := models.Target{}
	ta.HostName = "the beast"
	ta.Id = 1
	ta.Ipv4 = "129.211.213.123"
	ta.Ipv6 = "DEAD:BEEF::::"
	ta.CreatedAt = time.Now()
	ta.UpdatedAt = time.Now()
	cmdId := int64(1)
	response := []byte("some awesome response")
	ta.Path = fmt.Sprintf("/root/go/src/github.com/mojodojo101/c2server/internal_resources/targets/%v/", ta.Id)
	err = tu.SetCmdExecuted(ctx, &ta, cmdId, response)
	assert.NoError(t, err)

}
func TestFetchCmdResponse(t *testing.T) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	tr := targetdb.NewSQLRepo(db)
	cr := commanddb.NewSQLRepo(db)
	ctx := context.Background()
	cu := cmdusecase.NewCommandUsecase(cr, time.Second*2)
	tu := tusecase.NewTargetUsecase(tr, cu, time.Second*2)

	ta := models.Target{}
	ta.HostName = "the beast"
	ta.Id = 1
	ta.Ipv4 = "129.211.213.123"
	ta.Ipv6 = "DEAD:BEEF::::"
	ta.CreatedAt = time.Now()
	ta.UpdatedAt = time.Now()
	cmdId := int64(2)
	ta.Path = fmt.Sprintf("/root/go/src/github.com/mojodojo101/c2server/internal_resources/targets/%v/", ta.Id)
	buffer, err := tu.FetchCmdResponse(ctx, &ta, cmdId)
	assert.NoError(t, err)
	assert.NotEmpty(t, buffer)

}
