package store

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"github.com/zachlefevre/project_hopper_backend/com"
)

const (
	dbAddress = "event-db:26257"
)

type persistence interface {
	add(event *pb.Event) error
	getAll() []*pb.Event
	init()
}

type postgresDB struct {
}

func (db postgresDB) add(event *pb.Event) error {
	log.Printf("Adding event to DB", event)
	return nil
}

func (db postgresDB) getAll() []*pb.Event {
	log.Printf("request for all events")
	return nil
}

func (db postgresDB) init() {
	log.Printf("initializing DB")
	con, err := sql.Open("postgres", "postgresql://grace@"+dbAddress+"/?sslmode=disable")
	defer con.Close()
	if err != nil {
		log.Fatal("error connecting to the command DB: ", err)
	}

	if res, err := con.Exec(
		"CREATE DATABASE IF NOT EXISTS log"); err != nil {
		log.Fatal("cannot create database: ", err)
	} else {
		log.Printf("created database", res)
	}
	if res, err := con.Exec(
		"CREATE TABLE IF NOT EXISTS log.commands (id INT PRIMARY KEY, balance INT)"); err != nil {
		log.Fatal("cannot create table: ", err)
	} else {
		log.Printf("created table", res)
	}
}
