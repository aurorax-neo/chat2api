package router

import (
	"chat2api/app/conf"
	"chat2api/app/error_code"
	"chat2api/app/middleware"
	"chat2api/app/result"
	"chat2api/app/service"
	"chat2api/pkg/logx"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// routerRuntime 统一承载 API 路由模块初始化后的运行态变量。
// 这样做是为了避免模块级初始化变量散落，保持各模块初始化风格一致。
type routerRuntime struct {
	Engine *gin.Engine
	Server *http.Server
}

var runtime = &routerRuntime{}

const defaultHTTPShutdownTimeout = 10 * time.Second

func Init(ctx context.Context) func(ctx context.Context) {
	runtime.Engine = NewEngine()
	appConf := conf.GetApp()

	runtime.Server = &http.Server{
		Addr:    fmt.Sprintf("%v:%v", appConf.Bind, appConf.Port),
		Handler: runtime.Engine,
	}

	go func() {
		// 启动日志只记录绑定地址，不输出敏感配置，便于运维快速确认监听端口。
		logx.WithContext(ctx).Infof("httpServer started on http://%v:%v", appConf.Bind, appConf.Port)
		// ListenAndServe 为阻塞调用，放入协程避免卡住主线程信号监听流程。
		err := runtime.Server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			// 启动失败属于不可恢复错误：继续运行会导致服务无法对外提供 API，因此直接 Fatal 退出。
			logx.WithContext(ctx).Fatalf("httpServer start error: %v", err)
		}
	}()

	return func(ctx context.Context) {
		// 优雅停机：先停止接收新连接，再等待在途请求在超时窗口内收敛。
		// 为什么这样做：收到 SIGTERM/SIGINT 后直接退出会中断在途请求，可能导致调用方重试风暴或写入半途失败。
		// 输入与前置条件：cleanup ctx 来自应用退出阶段；若无 deadline，这里使用默认 10 秒收敛窗口。
		// 失败或异常行为：Shutdown 超时或返回错误会记录日志，但不会阻塞进程退出，避免卡死在清理阶段。
		logx.WithContext(ctx).Infof("http server is shutting down...")
		server := runtime.Server
		if server != nil {
			shutdownCtx, cancel := context.WithTimeout(ctx, defaultHTTPShutdownTimeout)
			defer cancel()
			if err := server.Shutdown(shutdownCtx); err != nil &&
				!errors.Is(err, http.ErrServerClosed) &&
				!errors.Is(err, context.Canceled) {
				logx.WithContext(ctx).Errorf("http server shutdown error: %v", err)
			}
		}
		runtime.Server = nil
		logx.WithContext(ctx).Infof("http server is shutting down complete")
		runtime.Engine = nil
	}
}

func NewEngine() *gin.Engine {
	// 强制使用 release 模式，避免 Gin 调试日志在生产/测试环境造成额外噪音和信息泄露风险。
	gin.SetMode(gin.ReleaseMode)
	// 采用 gin.New() 而非 gin.Default()，显式控制中间件装配，防止默认 Recovery/Logger 与项目日志规范冲突。
	engine := gin.New()

	// NoRoute 统一返回业务错误码，保证前端和调用方在“路径不存在”场景下拿到稳定的结构化响应。
	engine.NoRoute(func(ctx *gin.Context) {
		result.New(ctx).Error(error_code.NotFound)
	})

	// NoMethod 统一处理“路径存在但 HTTP 方法不支持”，避免返回框架默认文案导致客户端分支判断复杂。
	engine.NoMethod(func(ctx *gin.Context) {
		result.New(ctx).Error(error_code.MethodNotAllowed)
	})

	// 统一接入项目日志中间件，确保请求链路日志字段（trace/tag 等）与全局日志规范一致。
	engine.Use(middleware.Logger())

	// Register routes
	engine.GET("/", Index)
	engine.GET("/favicon.ico", Favicon)
	engine.GET("/favicon.png", Favicon)
	engine.GET("/ping", Ping)
	engine.GET("/admin", service.AdminConfigPage)
	adminRouter := engine.Group("/admin/api")
	adminRouter.Use(middleware.V1Cors)
	adminRouter.Use(middleware.AdminAuth)
	adminRouter.OPTIONS("/config", nil)
	adminRouter.GET("/config", service.GetAdminConfig)
	adminRouter.PUT("/config", service.SaveAdminConfig)
	apiRouter := engine.Group("/api")
	apiRouter.Use(middleware.V1Cors)
	apiRouter.OPTIONS("/usage/token", nil)
	apiRouter.GET("/usage/token", service.TokenUsage)
	apiRouter.OPTIONS("/usage/token/", nil)
	apiRouter.GET("/usage/token/", service.TokenUsage)
	v1Router := engine.Group("/v1")
	v1Router.Use(middleware.V1Cors)
	v1Router.Use(middleware.V1Auth)
	v1Router.GET("/accTokens", service.AccTokens)
	v1Router.OPTIONS("/usage", nil)
	v1Router.GET("/usage", service.Usage)
	v1Router.OPTIONS("/dashboard/billing/usage", nil)
	v1Router.GET("/dashboard/billing/usage", service.Usage)
	v1Router.OPTIONS("/dashboard/billing/subscription", nil)
	v1Router.GET("/dashboard/billing/subscription", service.BillingSubscription)
	v1Router.OPTIONS("/images/generations", nil)
	v1Router.POST("/images/generations", service.ImagesGenerations)
	v1Router.OPTIONS("/images/edits", nil)
	v1Router.POST("/images/edits", service.ImagesEdits)
	v1Router.OPTIONS("/chat/completions", nil)
	v1Router.POST("/chat/completions", service.Completions)
	v1Router.OPTIONS("/responses", nil)
	v1Router.POST("/responses", service.Responses)

	return engine
}

// HandlerGinEngine gin http handler
func HandlerGinEngine(w http.ResponseWriter, r *http.Request) {
	if runtime.Engine == nil {
		runtime.Engine = NewEngine()
	}
	runtime.Engine.ServeHTTP(w, r)
}

// Ping 测试接口
func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func Favicon(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Index 首页
func Index(c *gin.Context) {
	c.String(http.StatusOK, "hello, this is chat2api.")
}
