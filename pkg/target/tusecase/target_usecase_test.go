package tusecase_test

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/mojodojo101/c2server/pkg/command/commanddb"
	"github.com/mojodojo101/c2server/pkg/target/targetdb"

	"github.com/mojodojo101/c2server/pkg/target/usecase"
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
	ctx := context.Background()
	tu := usecase.NewTargetUsecase(tr, cr, time.Second*10)
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
	ctx := context.Background()
	tu := usecase.NewTargetUsecase(tr, cr, time.Second*10)

	ta := models.Target{}
	ta.CmdId = 1
	ta.HostName = "the beast"
	ta.Id = 1
	ta.Ipv4 = "129.211.213.123"
	ta.Ipv6 = "DEAD:BEEF::::"
	ta.CreatedAt = time.Now()
	ta.UpdatedAt = time.Now()
	ta.Path = fmt.Sprintf("/root/go/src/github.com/mojodojo101/c2server/internal_resources/targets/%v/", ta.Id)
	_, err = tu.Store(ctx, &ta)
	assert.NoError(t, err)

}
func TestUpdate(t *testing.T) {

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	tr := targetdb.NewSQLRepo(db)
	cr := commanddb.NewSQLRepo(db)
	ctx := context.Background()
	tu := usecase.NewTargetUsecase(tr, cr, time.Second*10)

	ta := models.Target{}
	ta.CmdId = 1
	ta.HostName = "the beast"
	ta.Id = 1
	ta.Ipv4 = "129.211.213.123"
	ta.Ipv6 = "DEAD:BEEF::::"
	ta.CreatedAt = time.Now()
	ta.UpdatedAt = time.Now()
	ta.Path = fmt.Sprintf("/root/go/src/github.com/mojodojo101/c2server/internal_resources/targets/%v/", ta.Id)
	err = tu.Update(ctx, &ta)
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
	tu := usecase.NewTargetUsecase(tr, cr, time.Second*10)

	ta := models.Target{}
	ta.CmdId = 1
	ta.HostName = "the beast"
	ta.Id = 1
	ta.Ipv4 = "129.211.213.123"
	ta.Ipv6 = "DEAD:BEEF::::"
	ta.CreatedAt = time.Now()
	ta.UpdatedAt = time.Now()
	ta.Path = fmt.Sprintf("/root/go/src/github.com/mojodojo101/c2server/internal_resources/targets/%v/", ta.Id)
	buffer, err := tu.FetchCmdResponse(ctx, &ta, ta.CmdId)
	assert.NoError(t, err)
	assert.NotEmpty(t, buffer)

}
