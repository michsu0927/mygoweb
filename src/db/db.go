package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

// Init initial Database connection
func Init() {
	// Db is database Object
	//db connection
	//var err error https://gorm.io/zh_CN/docs/gorm_config.html

	var dsn string
	var dialector gorm.Dialector

	dbEnv := os.Getenv("DB_ENV")
	if dbEnv == "sqlite" {
		dsn = os.Getenv("DATABASE_DSN")
		if dsn == "" {
			dsn = "data_db/test.db"
		}
		dialector = sqlite.Open(dsn)
	} else if dbEnv == "mysql" {
		dsn = os.Getenv("DATABASE_DSN")
		if dsn == "" {
			panic("DB_ENV=mysql but DATABASE_DSN is not set")
		}
		dialector = mysql.Open(dsn)
	} else {
		log.Printf("not found DB_ENV, will use sqlite as default")
		dsn = "data_db/test.db"
		dialector = sqlite.Open(dsn)

	}

	_, err := os.Stat("data_db")
	if os.IsNotExist(err) {
		os.Mkdir("data_db", 0755)
	}
	db, err = gorm.Open(dialector, &gorm.Config{PrepareStmt: true, SkipDefaultTransaction: true})
	if err != nil {
		if db != nil {
			fmt.Println(db.Error)
		}
		//panic("failed to connect database")
	}

	// now := time.Now()
	// //2025-02-12T07_58_11-log.txt
	// //formatted := fmt.Sprintf("%d-%02d-%02dT%02d_%02d_%02d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	// formatted := fmt.Sprintf("%d-%02d-%02d", now.Year(), now.Month(), now.Day())
	// f, err := os.OpenFile("log/"+formatted+"-db-log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Fatalf("error opening file: %v", err)
	// }
	// defer f.Close()

	// // show  sql string on console 這一段是把 sql 丟到 file
	// newLogger := logger.New(
	// 	log.New(f, "\r\n", log.LstdFlags), // io writer
	// 	logger.Config{
	// 		SlowThreshold: time.Second, // Slow SQL threshold
	// 		LogLevel:      logger.Info, // Log level
	// 		Colorful:      false,       // Disable color
	// 	},
	// )

	// db.Logger = newLogger
	// show sql string end

	db.AutoMigrate(&UserPointBalance{})
	db.AutoMigrate(&TransactionRecord{})
	db.AutoMigrate(&Task{})
}

// Manager Return Database Struct
func Manager() *gorm.DB {
	return db
}

// close function
// Close closes the database connection.
func Close() {
	dbEnv := os.Getenv("DB_ENV")
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Error getting database instance: %v", err)
	}
	if dbEnv == "sqlite" || dbEnv == "mysql" {
		if err := sqlDB.Close(); err != nil {
			log.Fatalf("Error closing database connection: %v", err)
		}
	} else {
		panic("DB_ENV is not sqlite or mysql, please setting the db ENV correctly")
	}
}

/*usage
Db := db.Manager()
rows, err := Db.Table("user").Select("player_id as playerID,password as Password,encrypt as Encrypt,status as Status,id as ID").Where("player_id = ?", un).Rows()
if err != nil {
	panic("failed to connect database")
}
defer rows.Close() //要記得 close
playerID := 0
password := ""
encrypt := ""
status := 0
ID := 0
for rows.Next() {
	err = rows.Scan(&playerID, &password, &encrypt, &status, &ID)
	if err != nil {
		fmt.Printf("Scan failed,err:%v\n", err)
		panic("Scan failed\n")
	}
}
*/
