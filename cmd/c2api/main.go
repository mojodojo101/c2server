package main


import(
	"database/sql"
	"context"
	"time"
	"fmt"

	_"github.com/lib/pq"

	"github.com/mojodojo101/c2server/pkg/client/database"
)


func main(){
	for i :=0;i<5;i++ {
		fmt.Println("")
	}
	fmt.Println("starting test")
	connStr := "host=localhost user=c2admin password=mojodojo101+ dbname=c2db port=5432 sslmode=require"
	db,err :=sql.Open("postgres",connStr)
	if err !=nil{
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	ctx :=context.Background()
	repo:=database.NewSQLClientRepo(db)
	err=repo.CreateClientTable(ctx)
	if err != nil {
		panic(err)
	}
	ip:="127.0.0.1"
	name:="mojo"
	password:="mojodojo101+"
	token:="myawesometoken"
	updated_at:=time.Now()
	created_at:=time.Now()
	client,err:=repo.CreateNewClient(ctx,ip,name,password,token,updated_at,created_at)
	fmt.Printf("%#v",client)
}