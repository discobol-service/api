package io


import (
	"sync"
	"github.com/jackc/pgx"
	"os"
	"context"
	"fmt"
)

type Pg struct {
	Connect *pgx.Conn
}

var once sync.Once
var pg *Pg

func GetPg() *Pg {
	once.Do(func() {
		pg = &Pg{connectToDarabase()}
	})
	return pg
}

func connectToDarabase() *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DISCOBOL_DBSTR"))
	if err != nil {
		panic( fmt.Sprintf("Unable to connect to database: %v\n", err.Error()) )
	}

	return conn
}