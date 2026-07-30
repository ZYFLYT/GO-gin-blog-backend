package main

import (
	"blog_project/internal/model"
	"blog_project/internal/router"
	"blog_project/pkg/conf"
	"blog_project/pkg/storage/orm"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/pflag"
)

var (
	cfgFile = pflag.StringP("config", "c", "./config/config.yml", "config file path")
)

func main() {
	// 1. 加载 .env 文件（必须放在最前面）
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ 未找到 .env 文件，将依赖系统环境变量")
	} else {
		log.Println("✅ 成功加载 .env 文件")
	}

	pflag.Parse()

	//2、加载配置
	cfg := conf.Init(*cfgFile)
	log.Printf("配置加载成功，当前模式：%s", cfg.App.Mode)

	// ========== 调试：打印关键环境变量 ==========
	log.Printf("🔍 BLOG_DSN = %s", os.Getenv("BLOG_DSN"))
	log.Printf("🔍 BLOG_ORM_DSN = %s", os.Getenv("BLOG_ORM_DSN"))
	// ==========================================

	log.Printf("✅ 配置加载成功，当前模式：%s", cfg.App.Mode)

	// ========== 调试：打印解析后的 DSN ==========
	log.Printf("🔍 cfg.ORM.Dsn = %s", cfg.ORM.Dsn)
	// ==========================================

	//3、连接数据库
	db := orm.NewMySQL(&orm.Config{
		Dsn:          cfg.ORM.Dsn,
		MaxIdleConns: cfg.ORM.MaxIdleConns,
		MaxOpenConns: cfg.ORM.MaxOpenConns,
	})

	//4、验证连接是否可用
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("获取底层sql.DB失败：%v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("数据库Ping失败：%v", err)
	}
	log.Println("数据库连接成功")

	//5、自动迁移，目前有user,category
	if err := db.AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.Tag{},
		&model.Article{},
		//后续添加 ...
	); err != nil {
		log.Fatalf("数据库迁移失败：%v", err)
	}
	log.Println("所有数据表创建/同步成功")

	//6、初始化路由
	r := router.InitRouter(db, cfg.JWT.Secret)

	//7、创建HTTP服务
	srv := &http.Server{
		Addr:    cfg.App.Addr,
		Handler: r,
	}

	// 8. 启动服务（协程）
	go func() {
		log.Printf("🚀 服务启动成功，监听地址: %s", cfg.App.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ 服务启动失败: %v", err)
		}
	}()

	// 9. 优雅关机
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("⏳ 正在关闭服务...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("❌ 强制关闭: %v", err)
	}
	log.Println("✅ 服务已安全关闭")
}
