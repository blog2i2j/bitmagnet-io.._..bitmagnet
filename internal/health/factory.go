package health

import (
	"github.com/bitmagnet-io/bitmagnet/internal/httpserver"
	"github.com/bitmagnet-io/bitmagnet/internal/lazy"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	Options []CheckerOption `group:"health_check_options"`
}

type Result struct {
	fx.Out
	Checker          lazy.Lazy[Checker]
	HTTPServerOption httpserver.Option `group:"http_server_options"`
}

func New(params Params) Result {
	lChecker := lazy.New[Checker](func() (Checker, error) {
		return NewChecker(params.Options...), nil
	})

	return Result{
		Checker:          lChecker,
		HTTPServerOption: handlerBuilder{lChecker},
	}
}

type handlerBuilder struct {
	Checker lazy.Lazy[Checker]
}

func (handlerBuilder) Key() string {
	return "health"
}

func (b handlerBuilder) Apply(e *gin.Engine) error {
	checker, err := b.Checker.Get()
	if err != nil {
		return err
	}

	handler := NewHandler(checker)

	e.GET("/status", func(c *gin.Context) {
		handler(c.Writer, c.Request)
	})

	return nil
}
