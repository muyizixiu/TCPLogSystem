package M

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type logger interface {
	log(map[string]interface{})
}
type mongo struct {
	url    string
	dbPool []db
}
type db struct {
	name string
	*mgo.Database
}
