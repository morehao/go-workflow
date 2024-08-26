package helper

import (
	"fmt"

	"github.com/morehao/go-tools/dbClient"
	"gorm.io/gorm"
)

var MysqlClient *gorm.DB

func InitDbClient() {
	mysqlClient, getMysqlClientErr := dbClient.InitMysql(Config.Mysql)
	if getMysqlClientErr != nil {
		panic(fmt.Sprintf("get mysql client error: %v", getMysqlClientErr))
	}
	MysqlClient = mysqlClient
}
