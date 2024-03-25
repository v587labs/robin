package log15

import (
	"github.com/inconshreveable/log15"
	"github.com/v587labs/robin/log"
)

type Log15 struct {
	log15.Logger
}

func (l *Log15) New(ctx ...interface{}) log.Logger {
	return &Log15{Logger: l.Logger.New(ctx...)}
}

func New(ctx ...interface{}) log.Logger {
	return &Log15{Logger: log15.New(ctx...)}
}
