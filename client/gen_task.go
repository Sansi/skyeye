package client

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

const dbconfig = "root:root@tcp(202.11.11.162:3306)/skyeye_db_test"

func gen_task() {
	db, err := sql.Open("mysql", dbconfig)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	stmtIns, err := db.Prepare("INSERT INTO `skyeye_db_test`.`task_todo`(`username`,`command_id`,`device_id`) VALUES ('qingpei', ?, 460020822485420);")
	if err != nil {
		panic(err.Error)
	}
	defer stmtIns.Close()

	for i := 1; i < 8; i++ {
		_, err = stmtIns.Exec(i)
		if err != nil {
			panic(err.Error())
		}
	}
}

func NewDevices(id []string) {
	db, err := sql.Open("mysql", dbconfig)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	fmt.Println(id)
	stmtIns, err := db.Prepare("INSERT INTO skyeye_db_test.device (id, cell, name, contract_no, device_type_id) VALUES (?,'12345678901','test_device','sh-2013-test',0);")
	if err != nil {
		panic(err.Error)
	}
	defer stmtIns.Close()

	for i := 0; i < len(id); i++ {
		_, err = stmtIns.Exec(id[i])
		if err != nil {
			panic(err.Error())
		}
	}
}

func DelDevices(id []string) {
	db, err := sql.Open("mysql", dbconfig)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	stmtIns, err := db.Prepare("DELETE FROM skyeye_db_test.device WHERE id =?;")
	if err != nil {
		panic(err.Error)
	}
	defer stmtIns.Close()

	for i := 0; i < len(id); i++ {
		_, err = stmtIns.Exec(id[i])
		if err != nil {
			panic(err.Error())
		}
	}
}
