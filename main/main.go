package main

import (
	"IdentityAuthentication/router"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始路由
	r := gin.Default()
	router.Router(r)

	r.Run()
}
