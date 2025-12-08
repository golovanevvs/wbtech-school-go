package transport

import (
	"log"
	"net/http"
	"net/http/pprof"
	"runtime"

	pkgPrometheus "analyzator/internal/pkg/pkgPrometheus"

	"github.com/gin-gonic/gin"
)

// NewRouter создает новый маршрутизатор
func NewRouter() *gin.Engine {
	router := gin.New()

	// Используем middleware для логирования
	router.Use(gin.Logger())

	// Используем middleware для восстановления после паники
	router.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		log.Printf("Паника: %v", recovered)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "внутренняя ошибка сервера",
		})
	}))

	// Добавляем middleware для сбора Prometheus метрик
	router.Use(pkgPrometheus.GinMiddleware())

	// Группа для API endpoints
	api := router.Group("/api/v1")

	// Health check endpoint
	api.GET("/health", healthHandler)

	// Memory и GC endpoints
	api.GET("/memory", memoryHandler)
	api.POST("/gc/trigger", triggerGCHandler)
	api.GET("/gc/status", gcStatusHandler)

	// Группа для Prometheus метрик
	metrics := router.Group("/metrics")
	{
		metrics.GET("", prometheusHandler)
	}

	// Группа для профилирования pprof
	pprofGroup := router.Group("/debug/pprof")
	{
		pprofGroup.GET("/", pprofIndexHandler)
		pprofGroup.GET("/cmdline", pprofCmdlineHandler)
		pprofGroup.GET("/profile", pprofProfileHandler)
		pprofGroup.GET("/symbol", pprofSymbolHandler)
		pprofGroup.GET("/trace", pprofTraceHandler)
		pprofGroup.GET("/goroutine", pprofGoroutineHandler)
		pprofGroup.GET("/heap", pprofHeapHandler)
		pprofGroup.GET("/allocs", pprofAllocsHandler)
		pprofGroup.GET("/block", pprofBlockHandler)
		pprofGroup.GET("/mutex", pprofMutexHandler)
	}

	return router
}

// healthHandler обработчик для проверки здоровья сервиса
func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"timestamp": gin.H{
			"unix": gin.H{
				"sec": gin.H{
					"int":   0,
					"int64": 0,
				},
			},
		},
		"version": "1.0.0",
	})
}

// memoryHandler обработчик для получения информации о памяти
func memoryHandler(c *gin.Context) {
	var memStats gin.H

	// Получаем статистику памяти через HTTP вызов к /metrics
	// В реальном приложении здесь была бы прямая интеграция
	resp, err := http.Get("http://localhost:" + c.Request.Host[len("http://"):] + "/metrics")
	if err != nil {
		log.Printf("Ошибка при получении метрик: %v", err)
		memStats = gin.H{
			"error": "не удалось получить статистику памяти",
		}
	} else {
		defer resp.Body.Close()
		// В реальном приложении здесь был бы парсинг ответа Prometheus
		memStats = gin.H{
			"message":          "используйте /metrics для получения подробной статистики памяти",
			"metrics_endpoint": "/metrics",
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"memory":      memStats,
		"description": "Получите подробную информацию о памяти и GC через /metrics",
	})
}

// triggerGCHandler обработчик для принудительного запуска GC
func triggerGCHandler(c *gin.Context) {
	// Запускаем принудительную сборку мусора
	runtime.GC()

	c.JSON(http.StatusOK, gin.H{
		"message": "Запущена принудительная сборка мусора",
		"timestamp": gin.H{
			"unix": gin.H{
				"sec": gin.H{
					"int":   0,
					"int64": 0,
				},
			},
		},
	})
}

// gcStatusHandler обработчик для получения статуса GC
func gcStatusHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"gc_status": gin.H{
			"enabled":     true,
			"description": "GC включен, используйте /metrics для получения детальной статистики",
		},
		"endpoints": gin.H{
			"metrics":    "/metrics - все метрики в формате Prometheus",
			"memory":     "/api/v1/memory - информация о памяти",
			"gc_trigger": "/api/v1/gc/trigger - принудительный запуск GC (POST)",
			"gc_status":  "/api/v1/gc/status - статус GC",
			"pprof":      "/debug/pprof/* - профилирование",
		},
	})
}

// prometheusHandler обработчик для метрик Prometheus
func prometheusHandler(c *gin.Context) {
	pkgPrometheus.Handler().ServeHTTP(c.Writer, c.Request)
}

// pprof handlers для профилирования
func pprofIndexHandler(c *gin.Context) {
	pprof.Index(c.Writer, c.Request)
}

func pprofCmdlineHandler(c *gin.Context) {
	pprof.Cmdline(c.Writer, c.Request)
}

func pprofProfileHandler(c *gin.Context) {
	pprof.Profile(c.Writer, c.Request)
}

func pprofSymbolHandler(c *gin.Context) {
	pprof.Symbol(c.Writer, c.Request)
}

func pprofTraceHandler(c *gin.Context) {
	pprof.Trace(c.Writer, c.Request)
}

func pprofGoroutineHandler(c *gin.Context) {
	pprof.Handler("goroutine").ServeHTTP(c.Writer, c.Request)
}

func pprofHeapHandler(c *gin.Context) {
	pprof.Handler("heap").ServeHTTP(c.Writer, c.Request)
}

func pprofAllocsHandler(c *gin.Context) {
	pprof.Handler("allocs").ServeHTTP(c.Writer, c.Request)
}

func pprofBlockHandler(c *gin.Context) {
	pprof.Handler("block").ServeHTTP(c.Writer, c.Request)
}

func pprofMutexHandler(c *gin.Context) {
	pprof.Handler("mutex").ServeHTTP(c.Writer, c.Request)
}
