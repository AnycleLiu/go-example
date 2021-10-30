package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

/*
create database user;
create table user.t_user(
	id int auto_increment primary key,
	userName varchar(150) not null,
	createdTime datetime not null default now()
)engine=innodb charset=utf8mb4;
*/

const dsn = "root:123456!@tcp(127.0.0.1:3306)/user?parseTime=true&loc=Local"

func main() {
	log.Println(time.Now().UTC())
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		log.Fatal("unable to use data source name", err)
	}

	defer db.Close()

	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(50)

	ping(db)

	printU990045 := func() {
		users, err := getuser(db, 990045)
		if err != nil {
			log.Fatal(err)
		}
		bs, err := json.Marshal(users)
		if err != nil {
			log.Fatal(err)
		}

		log.Println(string(bs))
	}

	printU990045()

	if err := updateUserCreatedTime(db, 990045, time.Now()); err != nil {
		log.Fatal(err)
	}
	log.Println("created time updated success.")

	printU990045()

}

func ping(db *sql.DB) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
}

func getuser(db *sql.DB, userid int) ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, "select id,userName,createdTime from t_user where id=?;", userid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*User, 0)
	for rows.Next() {
		u := &User{}
		if err := rows.Scan(&u.Id, &u.UserName, &u.CreatedTime); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	// If the database is being written to ensure to check for Close
	// errors that may be returned from the driver. The query may
	// encounter an auto-commit error and be forced to rollback changes.
	rerr := rows.Close()
	if rerr != nil {
		log.Fatal(rerr)
	}

	// Rows.Err will report the last error encountered by Rows.Scan.
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return users, nil
}

func updateUserCreatedTime(db *sql.DB, userid int, t time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, "update t_user set createdTime=? where id=?", t, userid)
	if err != nil {
		return err
	}

	return nil
}

type User struct {
	Id          int32
	UserName    string
	CreatedTime *time.Time `json:"createdTime"`
}
