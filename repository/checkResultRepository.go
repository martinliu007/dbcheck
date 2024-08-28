package repository

import (
	"database/sql"
	"dbcheck/config"
	"dbcheck/model"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

// SaveCheckResult 将检验结果写入数据库
func SaveCheckResult(checkResult []model.DBCheck,cfg *config.Config) error {
	// 连接数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", cfg.DBUsername, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("Failed to connect to database instance %s: %v", cfg.DBHost, err)
	}
	defer db.Close() // 确保关闭数据库连接

	// 清空校验结果表的历史数据
	_, err = db.Exec(`truncate table mysql_check`)
	if err != nil {
		log.Printf("Failed to truncate table mysql_check:%v", err)
		return err
	}

	// 向MySQL中插入巡检结果
	for _, row := range checkResult {
		_, err := db.Exec(`
            INSERT INTO mysql_check (db_type, ip_port, check_type, check_field, check_values) VALUES (?, ?, ?, ?, ?)`,
			row.DBType, row.IPPort, row.CheckType, row.CheckField, row.CheckValues)
		if err != nil {
			log.Printf("Failed to insert table mysql_check:%v", err)
		}
	}

	return nil
}
