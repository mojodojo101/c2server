package targetdb_test

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/mojodojo101/c2server/pkg/models"
	"github.com/mojodojo101/c2server/pkg/target/targetdb"
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
	br := targetdb.NewSQLRepo(db)
	err = br.CreateTable(ctx)
	assert.NoError(t, err)
	return
}

/*
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
	br := targetdb.NewSQLRepo(db)
	err = br.CreateTable(ctx)
	if err != nil {
		panic(err)
	}
	br.DropTable(ctx)
	assert.NoError(t, err)
	return
}
*/
func TestCreateNewTarget(t *testing.T) {
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
	br := targetdb.NewSQLRepo(db)
	err = br.CreateTable(ctx)
	if err != nil {
		panic(err)
	}
	ta := models.Target{}
	ta.Ipv4 = "192.168.138.145"
	ta.Ipv6 = "DEAD::::BEEF"
	ta.HostName = "mojo"
	ta.Path = fmt.Sprintf("$TARGETPATH/%v/", ta.Ipv4)
	ta.CreatedAt = time.Now()
	ta.UpdatedAt = time.Now()

	_, err = br.CreateNewTarget(ctx, &ta)
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
	br := targetdb.NewSQLRepo(db)
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

func TestUpdateCmdID(t *testing.T) {
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
	tr := targetdb.NewSQLRepo(db)
	err = tr.CreateTable(ctx)
	if err != nil {
		panic(err)
	}

	ta := models.Target{}
	ta.Ipv4 = "192.168.138.145"
	ta.Ipv6 = "DEAD::::BEEF"
	ta.HostName = "mojo"
	ta.Path = fmt.Sprintf("/root/go/src/github.com/mojodojo101/c2server/internal_resources/target/%v", ta.Ipv4)
	ta.UpdatedAt = time.Now()

	err = tr.Update(ctx, &ta)
	assert.NoError(t, err)
}

/*
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
	br := targetdb.NewSQLRepo(db)
	err = br.CreateTable(ctx)
	if err != nil {
		panic(err)
	}
	id := int64(1)
	err = br.DeleteByID(ctx, id)
	assert.NoError(t, err)
}
*/
