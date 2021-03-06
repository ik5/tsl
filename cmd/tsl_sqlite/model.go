// Copyright 2018 Yaacov Zamir <kobi.zamir@gmail.com>
// and other contributors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package main.
package main

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type book struct {
	ID     uint   `db:"id" json:"id"`
	Title  string `db:"title,omitempty" json:"title,omitempty"`
	Author string `db:"author,omitempty" json:"author,omitempty"`
	Pages  uint   `db:"pages,omitempty" json:"pages,omitempty"`
	Rating uint   `db:"Rating,omitempty" json:"Rating,omitempty"`
}

var books = []interface{}{
	book{Title: "Book", Author: "Joe", Pages: 100, Rating: 4},
	book{Title: "Other Book", Author: "Jane", Pages: 200, Rating: 3},
	book{Title: "Some Book", Author: "Jane", Pages: 50, Rating: 5},
	book{Title: "Some Other Book", Author: "Jane", Pages: 50},
	book{Title: "Good Book", Author: "Joe", Pages: 150, Rating: 4},
}

const sqlStmt = `
create table if not exists books (
	id integer not null primary key,
	title text,
	author text,
	pages integer,
	rating integer
);
delete from books;
`

func connect(ctx context.Context, url string) (tx *sql.Tx, err error) {
	var db *sql.DB

	db, err = sql.Open("sqlite3", url)
	check(err)

	tx, err = db.BeginTx(ctx, nil)

	return
}

func prepareCollection(ctx context.Context, tx *sql.Tx) (err error) {
	// Create table.
	fmt.Println("Createing table.")
	_, err = tx.ExecContext(ctx, sqlStmt)
	check(err)

	// Insert new books into the table.
	fmt.Println("Insert demo books.")
	stmt, err := tx.PrepareContext(ctx, "insert into books(title, author, pages, rating) values(?, ?, ?, ?)")
	check(err)

	defer stmt.Close()

	for _, b := range books {
		_, err = stmt.ExecContext(ctx,
			b.(book).Title,
			b.(book).Author,
			b.(book).Pages,
			b.(book).Rating)

		check(err)
	}

	return
}
