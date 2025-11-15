package handle

import (
	"github.com/gin-gonic/gin"
	"github.com/phuslu/log"
)

func InitRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	apiGroup := r.Group("/api")
	{
		dishGroup := apiGroup.Group("/dish")
		{
			dishGroup.GET("/list", ListDish)
		}
	}
	log.Info().Msg("router init success")
	return r
}
