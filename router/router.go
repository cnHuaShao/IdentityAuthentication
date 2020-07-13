package router

import (
	"IdentityAuthentication/cmd/common/model"
	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/google/uuid"
	"log"
	"net/http"
)

func routerHttp(r *gin.Engine) {
	// 加载登录页面
	r.LoadHTMLFiles("public/index.html")
	// 加载静态资源
	r.Static("/static", "static")
	// 设定基础登录页面
	r.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"titile": "身份认证系统",
		})
	})
}

// 设定认证api接口
func routerApi(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"msg": "pong",
			})
		})
	}
}

var OAuthClientDomain = "http://172.0.0.1"

func routerOAuth2(r *gin.Engine) {
	mg := manage.NewDefaultManager()
	// 设置默认令牌基础参数，详见源码中注释
	mg.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)
	// 1、存储token令牌
	mg.MustTokenStorage(store.NewMemoryTokenStore())

	// 2、客户端令牌存储
	// 创建客户端存储
	clientStore := store.NewClientStore()
	// 映射客户端存储接口
	mg.MapClientStorage(clientStore)

	// 3、创建一个默认的授权服务器
	srv := server.NewDefaultServer(mg)
	// 允许对令牌的GET请求
	srv.SetAllowGetAccessRequest(true)
	// 从请求中获取客户端信息，从表单中获取客户端数据 client_id  client_secret
	srv.SetClientInfoHandler(server.ClientFormHandler)

	// 内部错误处理
	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("内部错误：", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("响应错误：", re.Error.Error())
	})

	// 注册接口
	r.GET("/credentials", func(c *gin.Context) {
		// 生成一个客户端id
		clientId := uuid.New().String()[:8]
		// 生成一个客户端证书
		clientSecret := uuid.New().String()[:8]
		err := clientStore.Set(clientId, &models.Client{
			ID:     clientId,
			Secret: clientSecret,
			Domain: OAuthClientDomain,
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, model.Msg{
				"0x4000000",
				err.Error(),
				nil,
			})
			return
		}
		c.JSON(http.StatusOK, model.ReturnReqMess(model.MsgOK, map[string]string{
			"client_id":     clientId,
			"client_secret": clientSecret,
		}))
	})

	// 授权
	r.GET("/authorize", func(c *gin.Context) {
		err := srv.HandleAuthorizeRequest(c.Writer, c.Request)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.Msg{
				"0x400000",
				err.Error(),
				nil,
			})
		}
	})

	// 认证key获取
	r.GET("/token", func(c *gin.Context) {
		srv.HandleTokenRequest(c.Writer, c.Request)
	})
	r.Use(func(c *gin.Context) {
		_, err := srv.ValidationBearerToken(c.Request)
		if err != nil {
			c.JSON(http.StatusForbidden, model.Msg{
				"0x403000",
				err.Error(),
				nil,
			})
			c.Abort()
			return
		}
		c.Next()
	})
}

// 初始主路由
func Router(r *gin.Engine) {
	// 设定oauth2认证服务器
	routerOAuth2(r)
	// 设定静态资源与html
	routerHttp(r)
	// 设定认证api接口
	routerApi(r)
}
