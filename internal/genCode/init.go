package genCode

import (
	"go-workflow/internal/genCode/config"
	"os"
	"path/filepath"

	"github.com/morehao/go-tools/conf"
	"github.com/morehao/go-tools/dbClient"
	"github.com/morehao/go-tools/glog"
	"gorm.io/gorm"
)

var Config *config.Config
var MysqlClient *gorm.DB

func init() {
	// 初始化配置
	if workDir, err := os.Getwd(); err != nil {
		panic("get work dir error")
	} else {
		conf.SetAppRootDir(filepath.Join(workDir, "/internal/genCode"))
	}
	configFilepath := conf.GetAppRootDir() + "/config/config.yaml"
	var cfg config.Config
	conf.LoadConfig(configFilepath, &cfg)
	Config = &cfg

	// 初始化日志组件
	if err := glog.InitZapLogger(&Config.Log); err != nil {
		panic("glog initZapLogger error")
	}
	mysqlClient, getMysqlClientErr := dbClient.InitMysql(Config.Mysql)
	if getMysqlClientErr != nil {
		panic("get mysql client error")
	}
	MysqlClient = mysqlClient
}
