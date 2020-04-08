package clientdb_test

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/mojodojo101/c2server/pkg/client/clientdb"
	"github.com/mojodojo101/c2server/pkg/models"
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
	br := clientdb.NewSQLRepo(db)
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
	br := clientdb.NewSQLRepo(db)
	err = br.CreateTable(ctx)
	if err != nil {
		panic(err)
	}
	br.DropTable(ctx)
	assert.NoError(t, err)
	return
}
*/
func TestCreateNewClient(t *testing.T) {
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
	br := clientdb.NewSQLRepo(db)
	err = br.CreateTable(ctx)
	if err != nil {
		panic(err)
	}
	c := models.Client{}
	c.Ip = "127.0.0.1"
	c.Name = "mojo2"
	c.Password = "mojodojo101+"
	c.Token = "myawesometoken"
	c.CSRFToken = "myawesometokencsrf"
	c.UpdatedAt = time.Now()
	c.CreatedAt = time.Now()

	err = br.CreateNewClient(ctx, &c)
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
	br := clientdb.NewSQLRepo(db)
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

func TestUpdateHostName(t *testing.T) {
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
	cr := clientdb.NewSQLRepo(db)
	err = cr.CreateTable(ctx)
	if err != nil {
		panic(err)
	}
	id := int64(1)
	csrfToken := "32189038210938u90128093"
	err = cr.UpdateCSRFToken(ctx, id, csrfToken, time.Now())
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
	br := clientdb.NewSQLRepo(db)
	err = br.CreateTable(ctx)
	if err != nil {
		panic(err)
	}
	id := int64(1)
	err = br.DeleteByID(ctx, id)
	assert.NoError(t, err)
}
*/
