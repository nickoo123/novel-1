package app

import (
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/beego/i18n"

	"github.com/gomodule/redigo/redis"
	"github.com/vckai/novel/app/controllers"
	"github.com/vckai/novel/app/models"
	_ "github.com/vckai/novel/app/routers"
	"github.com/vckai/novel/app/services"
	"github.com/vckai/novel/app/utils"
	"github.com/vckai/novel/app/utils/log"
)

const (
	VERSION = "0.0.6"
)

var (
	RunTime   = time.Now()
	RedisPool *redis.Pool
)

func init() {
	// 初始化db
	models.InitDB()

	// 初始化redis
	initCache()

	// 初始化语言选项
	initLang()

	// 服务初始化配置
	services.Init()

	beego.AddFuncMap("urlfor", controllers.URLFor)

	// 注册模板函数
	utils.RegisterFuncMap()

	// 设置版本号
	beego.AppConfig.Set("version", VERSION)
}

func initCache() {
	var conn = beego.AppConfig.String("redisHost") + ":" + beego.AppConfig.String("redisPort")
	var password = beego.AppConfig.String("redisPass")
	var dbNum, _ = strconv.Atoi(beego.AppConfig.String("redisDB"))
	dialFunc := func() (c redis.Conn, err error) {
		c, err = redis.Dial("tcp", conn)
		if err != nil {
			return nil, err
		}

		if password != "" {
			if _, err := c.Do("AUTH", password); err != nil {
				c.Close()
				return nil, err
			}
		}

		_, selecterr := c.Do("SELECT", dbNum)
		if selecterr != nil {
			c.Close()
			return nil, selecterr
		}
		return
	}
	// initialize a new pool
	RedisPool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 180 * time.Second,
		Dial:        dialFunc,
	}
	InitPermissionRedisData()
}

//初始化缓存数据
func InitPermissionRedisData() {
	var redisPrefix = beego.AppConfig.String("redisPrefix")
	redisConn := RedisPool.Get()
	defer redisConn.Close()

	clearPermissionData(redisConn, redisPrefix)
}

//清除相关redis数据
func clearPermissionData(redisConn redis.Conn, redisPrefix string) {
	redisKeys, err := redis.Strings(redisConn.Do("KEYS", redisPrefix+"*"))
	if err != nil {
		return
	}
	redisConn.Do("MULTI")
	for _, v := range redisKeys {
		redisConn.Do("DEL", v)
	}
	redisConn.Do("EXEC")
}

//设置key value
func SetRedisKeyValue(key, value string) {
	var redisPrefix = beego.AppConfig.String("redisPrefix")
	redisConn := RedisPool.Get()
	defer redisConn.Close()
	redisConn.Do("SET", redisPrefix+key, value)
}

//删除key值
func DelRedisKeys(keys []string) {
	var redisPrefix = beego.AppConfig.String("redisPrefix")
	redisConn := RedisPool.Get()
	defer redisConn.Close()

	redisConn.Do("MULTI")
	for _, v := range keys {
		redisConn.Do("DEL", redisPrefix+v)
	}
	redisConn.Do("EXEC")
}

// 用key值获取数据
func GetRedisKeys(key string) interface{} {
	var redisPrefix = beego.AppConfig.String("redisPrefix")
	redisConn := RedisPool.Get()
	defer redisConn.Close()
	res, _ := redisConn.Do("GET", redisPrefix+key)
	return res
}

// 初始化语言选项
func initLang() {
	// Initialized language type list.
	langs := strings.Split(beego.AppConfig.String("lang::types"), "|")

	for _, lang := range langs {
		if err := i18n.SetMessage(lang, "lang/"+"locale_"+lang+".ini"); err != nil {
			log.Error("Fail to set message file: " + err.Error())
			return
		}
	}

	beego.AddFuncMap("i18n", i18n.Tr)
}
