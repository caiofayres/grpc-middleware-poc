package data

import (
	"errors"
	"sync"
)

type PersonDataService interface {
	Upsert(Person)
	Get(string) (*Person, error)
	Delete(string) error
}

type LocalPersonData struct {
	data sync.Map
}

type Person struct {
	Id string
	Name string
	Surname string
}

func NewPerson(id, name, surname string) Person {
	return Person{
		Id: id,
		Name: name,
		Surname: surname,
	}
}

var (
	ErrPersonNotFound error = errors.New("Person not found")
)

func NewLocalPersonData() PersonDataService {
	return &LocalPersonData{}
}

// inserts or update data
func (l *LocalPersonData)Upsert(p Person) {
	l.data.Store(p.Id, p)
}

func (l *LocalPersonData)Get(id string) (*Person, error) {
	v, ok := l.data.Load(id)
	if !ok {
		return nil, ErrPersonNotFound
	}
	p := v.(Person)
	return &p, nil
}

func (l *LocalPersonData)Delete(id string) error {
	_, loaded := l.data.LoadAndDelete(id)
	if !loaded {
		return ErrPersonNotFound
	}
	return nil
}