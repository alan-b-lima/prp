package auth

import "github.com/alan-b-lima/prp/pkg/uuid"

type Context struct {
	user  uuid.UUID
	level Level
}

func NewLogged(user uuid.UUID, level Level) Context {
	return Context{
		user:  user,
		level: level,
	}
}

func NewUnlogged() Context {
	return Context{level: Unlogged}
}

func (ctx *Context) User() uuid.UUID {
	return ctx.user
}

func (ctx *Context) Level() Level {
	return ctx.level
}
