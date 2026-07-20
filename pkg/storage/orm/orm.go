package orm

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// Config 对应 conf.ORMConfig 的字段（结构一致）
type Config struct {
	Dsn          string
	Driver       string
	Host         string
	Port         int
	User         string
	Password     string
	DBName       string
	Charset      string
	ParseTime    bool
	MaxIdleConns int
	MaxOpenConns int
}

func NewMySQL(c *Config) *gorm.DB {
	var dsn string
	if c.Dsn != "" {
		//优先使用完整的DSN(来自环境变量的BLOG_DSN)
		dsn = c.Dsn
	} else {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=Local",
			c.User, c.Password, c.Host, c.Port, c.DBName, c.Charset, c.ParseTime)
	}

	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	//打开数据库连接
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, //表名不加复数
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatalf("数据库连接失败：%v", err)
	}

	//获取底层sql.DB并设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("获取底层sqlDB失败：%v", err)
	}
	sqlDB.SetMaxIdleConns(c.MaxIdleConns) //空闲连接数
	sqlDB.SetMaxOpenConns(c.MaxOpenConns) //最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("数据库连接成功")
	return db
}
