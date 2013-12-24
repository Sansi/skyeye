package skyeye

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func gen_task() {
	db, err := sql.Open("mysql", "root:mariadb@sansi.com@tcp(localhost:3306)/skyeye_db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	stmtIns, err := db.Prepare("INSERT INTO `skyeye_db`.`task_todo`(`username`,`command_id`,`device_id`) VALUES ('qingpei', ?, 460020822485420);")
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
