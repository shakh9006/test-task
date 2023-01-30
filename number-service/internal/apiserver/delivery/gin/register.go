package delivery

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/shakh9006/numbers-store/internal/apiserver/middlewares"
)

type GinRouter struct {
	router     *gin.Engine
	numberCtrl *NumberCtrl
}

func NewGinRouter(numberCtrl *NumberCtrl) *GinRouter {
	return &GinRouter{
		router:     gin.New(),
		numberCtrl: numberCtrl,
	}
}

func (r *GinRouter) GetRouter() *gin.Engine {
	r.router.Use(gin.Recovery())

	r.router.SetTrustedProxies(nil)

	log.Info().Msg("Adding cors, request id and request logging middleware")
	r.router.Use(middlewares.CORSMiddleware())

	log.Info().Msg("Setting up routers")
	v1 := r.router.Group("/v1")
	{
		v1.GET("/number/:id", r.numberCtrl.GetById)
	}

	return r.router
}
