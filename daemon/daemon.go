package daemon

import (
	"github.com/gin-gonic/gin"
)

type Daemon struct {
	engine  *gin.Engine
	listen  string
}

func New(listen string) *Daemon {
	engine := gin.Default()
	Routes(engine)

	d := &Daemon{engine, listen}

	return d
}

func (d *Daemon) Run() {
	d.engine.Run(d.listen)
}