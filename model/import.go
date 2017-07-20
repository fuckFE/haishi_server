package model

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type OldType struct {
	ID   int
	Name string
}

type Forum struct {
	ID   int
	Name string
}

type Laws struct {
	Tname    string
	Fname    string
	Filename string
	Payload  string
}

func initData() {
	db, err := sql.Open("sqlite3", "./laws.db")
	if err != nil {
		panic(err)
	}

	ts, err := initType(db)
	if err != nil {
		panic(err)
	}
	fs, err := initForum(db)
	if err != nil {
		panic(err)
	}

	ls := make([]*Laws, 0)
	q, _ := db.Query("SELECT * FROM laws")
	for q.Next() {
		var id int
		var forumID string
		var forumName string
		var typeID string
		var typeName string
		var title string
		var contextHTML string
		var author string
		var remark string
		var contextText string
		var issueDate string
		var version string
		var deleted bool
		err := q.Scan(&id, &forumID, &forumName, &typeID, &typeName,
			&title, &contextHTML, &author, &remark, &contextText, &issueDate,
			&version, &deleted)
		if err != nil {
			panic(err)
		}

		l := &Laws{
			Tname:    typeName,
			Fname:    forumName,
			Filename: title,
			Payload:  contextHTML,
		}
		ls = append(ls, l)
	}

	newTypes := make([]*Type, 0)
	for _, t := range ts {
		newT, err := CreateType(t.Name, 1)
		if err != nil {
			panic(err)
		}
		newTypes = append(newTypes, newT)
		log.Println(newT)
	}

	for _, f := range fs {
		newF, err := CreateType(f.Name, 2)
		if err != nil {
			panic(err)
		}
		newTypes = append(newTypes, newF)
		log.Println(newF)
	}

	for _, l := range ls {
		b := &Book{
			Filename: l.Filename,
			Last:     time.Now(),
		}
		for _, t := range newTypes {
			if t.Name == l.Fname || t.Name == l.Tname {
				b.Types = append(b.Types, t.ID)
			}
		}

		if err := CreateBook(l.Filename, []byte(l.Payload), b.Types); err != nil {
			panic(err)
		}
	}

	log.Println("success")
}

func initType(db *sql.DB) ([]*OldType, error) {
	rows, err := db.Query("SELECT * FROM type")
	if err != nil {
		return nil, err
	}

	ts := make([]*OldType, 0)
	for rows.Next() {
		var id int
		var name string
		var exists bool
		if err := rows.Scan(&id, &name, &exists); err != nil {
			return nil, err
		}
		t := &OldType{
			ID:   id,
			Name: name,
		}
		ts = append(ts, t)
	}
	return ts, rows.Close()
}

func initForum(db *sql.DB) ([]*Forum, error) {
	rows, err := db.Query("SELECT * FROM forum")
	if err != nil {
		return nil, err
	}

	fs := make([]*Forum, 0)
	for rows.Next() {
		var id int
		var name string
		var exists bool
		if err := rows.Scan(&id, &name, &exists); err != nil {
			return nil, err
		}
		f := &Forum{
			ID:   id,
			Name: name,
		}
		fs = append(fs, f)
	}
	return fs, rows.Close()
}
