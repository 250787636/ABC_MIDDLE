package automigrate

import (
	"errors"
	"example.com/m/model"
	"example.com/m/pkg/log"
	"fmt"
	"gorm.io/gorm"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"
)

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}
type TableVersion struct {
	Struct  interface{}
	Version string
}
type AutoMigrate struct {
	DB *gorm.DB
	DBConfig
	DataTables []TableVersion
	ToolTables []TableVersion // 固定数据表
	SqlPathMap map[string]string
}

// 获取所有sql文件的文件名（目录下面的sql文件不能重名）
func NewAutoMigrate(db *gorm.DB, sqlFilesDir string, dbconfig DBConfig) (migrate *AutoMigrate) {
	migrate = new(AutoMigrate)
	migrate.DBConfig = dbconfig
	migrate.DB = db
	migrate.SqlPathMap = make(map[string]string)
	match, err := filepath.Glob(sqlFilesDir + "/*.sql")
	if err != nil {
		log.Error(err)
	}
	for _, filePath := range match {
		fileName := filepath.Base(filePath)
		ext := filepath.Ext(fileName)
		tableName := fileName[0 : len(fileName)-len(ext)]
		migrate.SqlPathMap[tableName] = filePath
	}
	return migrate
}

func (t *AutoMigrate) AutoMigrate() {
	migrator := t.DB.Migrator()
	if err := migrator.AutoMigrate(model.GlobalVariable{}); err != nil {
		log.Error(err)
	}
	// 读取sql文件 进行赋值
	for _, table := range t.DataTables {
		_, tableName := t.GetStructAndTableName(table.Struct)
		if !migrator.HasTable(tableName) {
			if sqlFilePath, ok := t.SqlPathMap[tableName]; ok {
				log.Infof("开始初始化表%s", tableName)
				if err := t.ExecSqlFile(sqlFilePath); err != nil {
					log.Errorf("初始化表%s失败:%s", tableName, err.Error())
				} else {
					log.Infof("初始化表%s完成", tableName)
				}
			}
			if err := t.DB.AutoMigrate(table.Struct); err != nil {
				log.Errorf("表%s:%v", table, err)
			}
			if err := t.UpdateTableVersion(tableName, table.Version); err != nil {
				log.Errorf("表%s:%v", table, err)
			}
		} else {
			t.UpdateTable(table)
		}
	}
	// 工具表，直接执行sql文件删表重建，若表中有运行时数据，需做sql文件执行前保存数据并在sql文件执行后更新
	for _, table := range t.ToolTables {
		_, tableName := t.GetStructAndTableName(table.Struct)
		if !migrator.HasTable(tableName) {
			log.Infof("开始初始化表%s", tableName)
			if err := t.ExecSqlFile(t.SqlPathMap[tableName]); err != nil {
				log.Errorf("初始化表%s失败", tableName, err.Error())
			} else {
				log.Infof("初始化表%s完成", tableName)
				if err = t.UpdateTableVersion(tableName, table.Version); err != nil {
					log.Errorf("表%s:%v", tableName, err)
				}
			}
		} else {
			t.UpdateTable(table)
		}
	}
}

// 更新表
func (t *AutoMigrate) UpdateTable(table TableVersion) {
	structName, tableName := t.GetStructAndTableName(table.Struct)
	newVersion, currentVersion := table.Version, t.TableVersion(tableName)
	if currentVersion == newVersion {
		return
	}
	log.Infof("升级表%s:%s -> %s", tableName, currentVersion, newVersion)
	if err := t.DB.AutoMigrate(table.Struct); err != nil {
		log.Error("表%s:%v", tableName, err)
	}
	var err error
	if strings.HasPrefix(newVersion, "#") {
		err = t.UpgradeTable(structName, currentVersion, newVersion)
	} else {
		err = t.UpdateTableVersion(tableName, newVersion)
	}
	if err != nil {
		log.Errorf("表%s升级失败:%v", tableName, err)
	} else {
		log.Infof("表%s升级完成", tableName)
	}
}

// 通过反射调用方法
func (t *AutoMigrate) UpgradeTable(structName, currentVersion, newVersion string) error {
	methodName := "Upgrade" + structName
	tableName := t.DB.NamingStrategy.TableName(structName)
	method := reflect.ValueOf(t).MethodByName(methodName)
	if method.Kind() != reflect.Invalid {
		upgradeFunc, ok := method.Interface().(func(string, string, string) error)
		if !ok {
			return errors.New(fmt.Sprintf("方法%s签名错判,真缺为:func(string,string,string)error", methodName))
		} else {
			if err := upgradeFunc(tableName, currentVersion, newVersion); err != nil {
				return err
			}
			if err := t.UpdateTableVersion(tableName, newVersion); err != nil {
				return err
			}
		}
	} else {
		return errors.New(fmt.Sprintf("未实现%s方法,无法升级表%s", methodName, tableName))
	}
	return nil
}

// 获取结构体名 和 表名
func (t *AutoMigrate) GetStructAndTableName(tableStruct interface{}) (string, string) {
	structName := reflect.TypeOf(tableStruct).Name()
	tableName := t.DB.NamingStrategy.TableName(structName)
	return structName, tableName
}

// 获取 数据库中的表版本号
func (t *AutoMigrate) TableVersion(tableName string) (currentVersion string) {
	var variable model.GlobalVariable
	if err := t.DB.Where("type = ? AND name = ?",
		"table_version", tableName).First(&variable).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error(err)
		return
	}
	return variable.Value
}

func (t *AutoMigrate) UpdateTableVersion(tableName, newVersion string) error {
	typ := "table_version"
	var variable model.GlobalVariable
	if err := t.DB.Where("type = ? and name = ?", typ,
		tableName).First(&variable).Error; err == nil {
		if err = t.DB.Model(&model.GlobalVariable{}).Where("type = ? and name = ?",
			typ, tableName).Update("value", newVersion).Error; err != nil {
			return err
		}
	} else if err == gorm.ErrRecordNotFound {
		variable = model.GlobalVariable{
			Type:  typ,
			Name:  tableName,
			Value: newVersion,
		}
		if err = t.DB.Create(&variable).Error; err != nil {
			return err
		}
	}
	return nil
}

func (t *AutoMigrate) ExecSqlFile(sqlPath string) error {
	fileBytes, err := ioutil.ReadFile(sqlPath)
	if err != nil {
		return err
	}
	if err = t.DB.Exec(string(fileBytes)).Error; err != nil {
		return err
	}
	return nil
}
