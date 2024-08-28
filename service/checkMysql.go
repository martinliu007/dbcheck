package service

import (
	"database/sql"
	"dbcheck/config"
	"dbcheck/model"
	"dbcheck/repository"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func CheckMysql(cfg *config.Config) error {
	var checkResult []model.DBCheck

	for _, mysqlHost := range cfg.CheckDBList {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/information_schema", cfg.CheckDBUsername, cfg.CheckDBPassword, mysqlHost, cfg.CheckDBPort)
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			log.Printf("Failed to connect to database instance %s: %v", mysqlHost, err)
			continue // 继续尝试下一个数据库实例
		}
		defer db.Close() // 确保关闭数据库连接

		// 巡检表行数top10
		rowsQuery := "select concat(table_schema,'.',table_name), table_rows from information_schema.tables where table_schema not in " +
			"('information_schema', 'mysql', 'performance_schema', 'sys') order by table_rows desc limit 10"

		tableRowsInfo, err := GetCheckResult(db, mysqlHost,rowsQuery,"TableRows",cfg)
		if err != nil {
			log.Printf("Failed to check tables for database instance %s: %v", mysqlHost, err)
			continue // 继续尝试下一个数据库实例
		}
		checkResult = append(checkResult, tableRowsInfo...)

		// 巡检表存储引擎是非InnoDB的
		engineQuery := " select concat(table_schema,'.',table_name), engine from information_schema.tables where table_schema not in " +
			"('information_schema', 'mysql', 'performance_schema', 'sys') and engine <> 'innodb'"
		tableEngineInfo, err := GetCheckResult(db, mysqlHost,engineQuery,"TableEngine",cfg)
		if err != nil {
			log.Printf("Failed to check tables for database instance %s: %v", mysqlHost, err)
			continue // 继续尝试下一个数据库实例
		}
		checkResult = append(checkResult, tableEngineInfo...)

		// 巡检自增值使用率top10
		autoIncrementQuery := "select concat(table_schema,'.',table_name), round((auto_increment / pow(2, 31)) * 100, 2) as 'auto_increment_usage' " +
			"from information_schema.tables where table_schema NOT IN ('information_schema', 'mysql', 'performance_schema', 'sys') " +
			"and auto_increment is not null order by auto_increment_usage desc limit 10"
		tableAutoIncrementInfo, err := GetCheckResult(db, mysqlHost,autoIncrementQuery,"TableAutoIncrement",cfg)
		if err != nil {
			log.Printf("Failed to check tables for database instance %s: %v", mysqlHost, err)
			continue // 继续尝试下一个数据库实例
		}
		checkResult = append(checkResult, tableAutoIncrementInfo...)

		// 巡检表碎片率top10
		fragmentationQuery := "select concat(table_schema,'.',table_name),if(data_length > 0, round(data_free/(data_length + index_length + data_free)" +
			" * 100, 2), 0) as 'fragmentation' from information_schema.tables  where table_schema not in ('information_schema', 'mysql', " +
			"'performance_schema', 'sys') and engine = 'innodb' order by fragmentation desc limit 10"
		tableFragmentationInfo, err := GetCheckResult(db, mysqlHost,fragmentationQuery,"TableFragmentation",cfg)
		if err != nil {
			log.Printf("Failed to check tables for database instance %s: %v", mysqlHost, err)
			continue // 继续尝试下一个数据库实例
		}
		checkResult = append(checkResult, tableFragmentationInfo...)


		// 巡检重要参数
		variablesQuery := "show global variables where Variable_name in ('innodb_flush_log_at_trx_commit','sync_binlog'," +
			"'binlog_format','character_set_server','system_time_zone');"
		variablesInfo, err := GetCheckResult(db, mysqlHost,variablesQuery,"Variables",cfg)
		if err != nil {
			log.Printf("Failed to check tables for database instance %s: %v", mysqlHost, err)
			continue // 继续尝试下一个数据库实例
		}
		checkResult = append(checkResult, variablesInfo...)


	}
	for _, info := range checkResult {
		fmt.Println(info)
	}

	err := repository.SaveCheckResult(checkResult,cfg)
	if err != nil {
		log.Printf("Failed to save check result: %v" , err)
		 // 继续尝试下一个数据库实例
	}

	return nil // 所有数据库实例处理完成
}