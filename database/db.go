package database

import (
	"log"

	"github.com/jmoiron/sqlx"
)

type database struct {
	DB *sqlx.DB
}

func InitDb(uri string) (*database, error) {
	db, err := sqlx.Connect("mysql", uri)
	if err != nil {
		log.Fatalln(err)
	}

	var schema = []string{
		`CREATE DATABASE IF NOT EXISTS tickets;`,
		`USE tickets;`,
		`CREATE TABLE IF NOT EXISTS users (
			id int(11) unsigned NOT NULL AUTO_INCREMENT,
			email varchar(32) DEFAULT NULL,
			password varchar(128) DEFAULT NULL,
			api_key text DEFAULT NULL,
			PRIMARY KEY (id)
		  )`,
		`CREATE TABLE IF NOT EXISTS events (
			id int(11) unsigned NOT NULL AUTO_INCREMENT,
			event_name varchar(128) DEFAULT NULL,
			address varchar(128) DEFAULT NULL,
			handphone varchar(13) DEFAULT NULL,
			bank_number varchar(64) DEFAULT NULL,
			id_user int(11) unsigned NOT NULL,
			PRIMARY KEY (id),
			FOREIGN KEY fk_user(id_user)
			REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS products (
			id int(11) unsigned NOT NULL AUTO_INCREMENT,
			ticket_name varchar(128) DEFAULT NULL,
			ticket_class varchar(128) DEFAULT NULL,
			quantity int(11) DEFAULT NULL,
			description text DEFAULT NULL,
			price int(11) DEFAULT NULL,
			url_pic varchar(128) DEFAULT NULL,
			id_event int(11) unsigned NOT NULL,
			PRIMARY KEY (id),
			FOREIGN KEY fk_event(id_event)
			REFERENCES events(id)
		)`,
		`CREATE TABLE IF NOT EXISTS customers (
			id int(11) unsigned NOT NULL AUTO_INCREMENT,
			cust_name int(11) DEFAULT NULL,
			cust_email varchar(128) NOT NULL,
			cust_address int(11) DEFAULT NULL,
			id_event int(11) unsigned NOT NULL,
			PRIMARY KEY (id),
			FOREIGN KEY fk_event2(id_event)
			REFERENCES events(id)
		)`,
		`CREATE TABLE IF NOT EXISTS bookings (
			id int(11) unsigned NOT NULL AUTO_INCREMENT,
			id_ticket int(11) unsigned NOT NULL,
			id_customer int(11) unsigned NOT NULL,
			quantity int(11) DEFAULT NULL,
			total int(11) DEFAULT NULL,
			status enum('1','2') DEFAULT NULL,
			PRIMARY KEY (id),
			FOREIGN KEY fk_ticket(id_ticket)
			REFERENCES products(id),
			FOREIGN KEY fk_customer(id_customer)
			REFERENCES customers(id)
		)`,
	}

	for _, value := range schema {
		db.MustExec(value)
	}
	return &database{DB: db}, nil
}

func (d *database) GetDB() *sqlx.DB {
	return d.DB
}
