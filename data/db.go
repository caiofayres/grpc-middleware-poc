package data

import "github.com/jmoiron/sqlx"

type PersonDatabase struct {
	db *sqlx.DB
}

func NewPersonDatabase(db *sqlx.DB) *PersonDatabase {
	return &PersonDatabase{db: db}
}

func (p *PersonDatabase) Upsert(person Person) error {
	_, err := p.db.NamedExec("INSERT INTO person (id, name, surname) VALUES (:id, :name, :surname) ON CONFLICT (id) DO UPDATE SET name = :name, surname = :surname", person)
	return err
}

func (p *PersonDatabase) Get(id string) (*Person, error) {
	var person Person
	err := p.db.Get(&person, "SELECT * FROM person WHERE id = $1", id)
	return &person, err
}

func (p *PersonDatabase) Delete(id string) error {
	_, err := p.db.Exec("DELETE FROM person WHERE id = $1", id)
	return err
}

func ConnectToDatabase() (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", "user=postgres dbname=postgres password=postgres sslmode=disable")
	if err != nil {
		return nil, err
	}
	return db, nil
}
