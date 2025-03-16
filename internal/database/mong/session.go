package mong

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

// ISession is interface of session structure.
type ISession interface {
	// DoTransaction is process session function with transaction.
	DoTransaction(fn SessionFunc) (any, error)

	// Close for end of session process.
	Close()
}

// Session is mongo.Session structure wrapper.
type Session struct {
	ctx     context.Context
	session mongo.Session
}

// SessionContext is mongo.SessionContext wrapper.
type SessionContext mongo.SessionContext

// SessionFunc is template of session callback function.
type SessionFunc func(SessionContext) (any, error)

func (sess Session) DoTransaction(fn SessionFunc) (any, error) {
	callback := func(sctx mongo.SessionContext) (any, error) {
		return fn(sctx)
	}

	return sess.session.WithTransaction(sess.ctx, callback)
}

func (sess Session) Close() {
	sess.session.EndSession(sess.ctx)
}
