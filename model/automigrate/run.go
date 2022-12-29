package automigrate

import "example.com/m/pkg/app"

func Run() {
	mysqlConf := app.Conf.Mysql
	migrate := NewAutoMigrate(app.DB, app.ResourceDir+"/sql", DBConfig{
		Host:     mysqlConf.Host,
		Port:     mysqlConf.Port,
		Username: mysqlConf.Username,
		Password: mysqlConf.Password,
		Database: mysqlConf.Database,
	})
	setTable(migrate)
	migrate.AutoMigrate()
}
