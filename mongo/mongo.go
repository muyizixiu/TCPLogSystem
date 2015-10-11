package M

import (
	"gopkg.in/mgo.v2"
	"time"
)

var (
	M *mongo
)

type logger interface {
	Log(map[string]interface{}) bool
}
type mongo struct {
	url      string
	dbPool   map[string]*db
	session  *mgo.Session
	maxTries int
}

func (m mongo) Log(v map[string]interface{}) bool {
	dbname, ok := v["__db"]
	if !ok || dbname == nil {
		return false
	}
	name, ok0 := dbname.(string)
	if !ok0 {
		return false
	}
	delete(v, "__db")
	db := m.getDB(name)
	cname, ok1 := v["__c"]
	if _, ok := cname.(string); !ok {
		return false
	}
	if !ok1 {
		return false
	}
	delete(v, "__c")
	c := db.DB.C(cname.(string))
	err := c.Insert(&v)
	if err != nil {
		return false
	}
	return true
}
func (m *mongo) getDB(name string) *db {
	if Db, ok := m.dbPool[name]; ok {
		if Db != nil {
			return Db
		}
	}
	m.dbPool[name] = newDB(name, m.session.DB(name))
	return m.dbPool[name]

}
func init() {
	M = &mongo{
		url: "localhost:27017",
	}
}
func (m *mongo) open() bool {
	if m.session != nil {
		if m.session.Ping() == nil {
			return true
		}
	}
	var err error
	m.session, err = mgo.Dial(m.url)
	if err != nil {
		return false
	}
	return true
}

//最大次数尝试连接
func (m *mongo) openMaxTries() bool {
	for i := 0; i < m.maxTries; {
		i++
		if m.open() {
			return true
		}
	}
	return false
}

type db struct {
	name         string
	DB           *mgo.Database
	lastUsedTime time.Time
	count        int
}

func newDB(n string, Db *mgo.Database) *db {
	return &db{name: n, DB: Db, lastUsedTime: time.Now()}
}
