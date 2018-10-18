package db

import (
	"log"
	"sync"
	"time"

	"github.com/das-frama/website/app"
	"github.com/das-frama/website/app/session"

	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

var provider = &Provider{}

type Provider struct {
	lock sync.Mutex
}

type SessionStore struct {
	ID         string                      `rethinkdb:"id,omitempty"`
	SID        string                      `rethinkdb:"sid"`
	AccessedAt time.Time                   `rethinkdb:"accessed_at"`
	Value      map[interface{}]interface{} `rethinkdb:"value"`
}

func init() {
	session.Register("db", provider)
}

func (st *SessionStore) Set(key, value interface{}) error {
	st.Value[key] = value
	_, err := r.Table("session").GetAllByIndex("sid", st.SID).Update(map[string]interface{}{
		"value":       st.Value,
		"accessed_at": time.Now(),
	}).RunWrite(app.RethinkSession)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (st *SessionStore) Get(key interface{}) interface{} {
	// Обновить сессию.
	_, err := r.Table("session").GetAllByIndex("sid", st.SID).Update(map[string]interface{}{
		"accessed_at": time.Now(),
	}).RunWrite(app.RethinkSession)
	if err != nil {
		log.Println(err)
		return nil
	}

	if v, ok := st.Value[key]; ok {
		return v
	}

	return nil
}

func (st *SessionStore) Delete(key interface{}) error {
	delete(st.Value, key)
	_, err := r.Table("session").GetAllByIndex("sid", st.SID).Update(map[string]interface{}{
		"value":       st.Value,
		"accessed_at": time.Now(),
	}).RunWrite(app.RethinkSession)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (st *SessionStore) SessionID() string {
	return st.SID
}

func (p *Provider) SessionInit(sid string) (session.Session, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	val := make(map[interface{}]interface{}, 0)
	session := &SessionStore{
		SID:        sid,
		AccessedAt: time.Now(),
		Value:      val,
	}
	_, err := r.Table("session").Insert(session).RunWrite(app.RethinkSession)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return session, nil
}

func (p *Provider) SessionRead(sid string) (session.Session, error) {
	res, err := r.Table("session").GetAllByIndex("sid", sid).Run(app.RethinkSession)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if res.IsNil() {
		session, err := p.SessionInit(sid)
		return session, err
	} else {
		session := &SessionStore{}
		err = res.One(&session)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return session, nil
	}
}

func (p *Provider) SessionDestroy(sid string) error {
	_, err := r.Table("session").GetAllByIndex("sid", sid).Delete().RunWrite(app.RethinkSession)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (p *Provider) SessionGC(maxlifetime int64) {
	p.lock.Lock()
	defer p.lock.Unlock()

	_, err := r.Table("session").Filter(
		r.Row.Field("accessed_at").Lt(time.Now().Unix() - maxlifetime),
	).Delete().RunWrite(app.RethinkSession)
	if err != nil {
		log.Println(err)
	}
}
