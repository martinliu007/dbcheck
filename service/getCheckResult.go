package service

import (
	"database/sql"
	"dbcheck/config"
	"dbcheck/model"
	"log"
)

func GetCheckResult(db *sql.DB, mysqlHost,query,checkType string, cfg *config.Config) ([]model.DBCheck, error) {
	var tableRowsInfo []model.DBCheck

	// 执行巡检逻辑
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Failed to query source database: %v", err)
		return tableRowsInfo, err
	}
	defer rows.Close() // 确保关闭rows

	for rows.Next() {
		var ti model.DBCheck

		if err := rows.Scan(&ti.CheckField, &ti.CheckValues); err != nil {
			log.Printf("Failed to scan row: %v", err)
			return tableRowsInfo, err
		}

		ti.DBType = "MySQL"
		ti.IPPort = mysqlHost + ":" + cfg.CheckDBPort
		ti.CheckType = checkType
		tableRowsInfo = append(tableRowsInfo, ti)
	}

	return tableRowsInfo, nil
}