package cusecase_test

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mojodojo101/c2server/pkg/client/clientdb"
	"github.com/mojodojo101/c2server/pkg/client/cusecase"
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

/*
func TestCreateTable(t *testing.T) {

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	cr := clientdb.NewSQLRepo(db)
	ctx := context.Background()
	cu := cusecase.NewClientUsecase(cr, time.Second*10)
	err = cu.CreateTable(ctx)
	assert.NoError(t, err)

}

func TestStore(t *testing.T) {

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	cr := clientdb.NewSQLRepo(db)
	ctx := context.Background()
	cu := cusecase.NewClientUsecase(cr, time.Second*10)
	c := models.Client{}
	c.Id = 1
	c.Password = "mypass"
	c.Token = "ndwsioanwoqinoid"
	c.CSRFToken = "32132132121"
	c.Ip = "192.168.122.102"
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	err = cu.Store(ctx, &c)
	assert.NoError(t, err)

}
*/

func TestListTargetCommands(t *testing.T) {

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	timeout := time.Second * 5
	cr := clientdb.NewSQLRepo(db)
	cmdr := commanddb.NewSQLRepo(db)
	tr := targetdb.NewSQLRepo(db)
	cmdu := cmdusecase.NewCommandUsecase(cmdr, timeout)
	tu := tusecase.NewTargetUsecase(tr, cmdu, timeout)
	ctx := context.Background()
	cu := cusecase.NewClientUsecase(cr, tu, timeout)
	c := models.Client{}
	c.Id = 1
	c.Name = "mojo"
	c.Password = "mojodojo101+"
	c.Token = "this should be different"
	c.CSRFToken = "32132132121"
	c.Ip = "192.168.122.102"
	c.UpdatedAt = time.Now()
	c.CreatedAt = time.Time{}
	tId := int64(1)
	amount := int64(20)
	cmds, err := cu.ListTargetCommands(ctx, &c, tId, amount)
	fmt.Printf("\nCmds=%v\n", cmds)
	assert.NoError(t, err)
}

/*
func TestUpdate(t *testing.T) {

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	cr := clientdb.NewSQLRepo(db)
	ctx := context.Background()
	cu := cusecase.NewClientUsecase(cr, time.Second*10)
	c := models.Client{}
	c.Id = 1
	c.Password = "mypass"
	c.Token = "this should be different"
	c.CSRFToken = "32132132121"
	c.Ip = "192.168.122.102"
	c.UpdatedAt = time.Now()
	c.CreatedAt = time.Time{}

	err = cu.Update(ctx, &c)
	assert.NoError(t, err)

}
*/
