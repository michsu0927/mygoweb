package controller

import (
	"bytes"

	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"web/src/db"
	"web/src/lib"

	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
)

func AddTask(c echo.Context) error {

	// var json map[string]interface{} = map[string]interface{}{}
	// if err := c.Bind(&json); err != nil {
	// 	return err
	// }
	// return c.String(http.StatusOK, fmt.Sprintf("%v", json))
	// Get post json contents
	var jsonData map[string]interface{} = map[string]interface{}{}
	err := c.Bind(&jsonData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, lib.PyDict{
			"success": false,
			"message": "Invalid JSON payload",
		})
	}

	// log the contents of c.Request().Body to a log file named log/yyyy-mm-dd.log
	// bodyContent, _ := io.ReadAll(c.Request().Body)
	lib.Log(string(fmt.Sprintf("%v", jsonData)))
	lib.Log(string(fmt.Sprintf("%v", jsonData["TaskID"])))
	//insert in to task
	//check not empty
	if taskID, ok := jsonData["TaskID"]; !ok || strings.TrimSpace(string(fmt.Sprintf("%v", taskID))) == "" {
		return c.JSON(http.StatusBadRequest, lib.PyDict{
			"success": false,
			"message": "TaskID cannot be empty",
		})
	}
	if userID, ok := jsonData["UserID"]; !ok || strings.TrimSpace(string(fmt.Sprintf("%v", userID))) == "" {
		return c.JSON(http.StatusBadRequest, lib.PyDict{
			"success": false,
			"message": "UserID cannot be empty",
		})
	}

	task := db.Task{
		TaskID:          "task1",
		UserID:          "user1",
		TaskType:        "type1",
		Description:     "desc1",
		PointsChange:    10,
		ExpiredDatetime: nil,
		//CreateDatetime:  time.Now(),
		Status: 0, // Set the default value for status
	}

	task.TaskID = strings.ToLower(string(fmt.Sprintf("%v", jsonData["TaskID"])))
	if userID, ok := jsonData["UserID"]; ok {
		task.UserID = strings.ToLower(string(fmt.Sprintf("%v", userID)))
	}

	if taskType, ok := jsonData["TaskType"]; ok {
		task.TaskType = strings.ToLower(string(fmt.Sprintf("%v", taskType)))
	} else {
		task.TaskType = "" // Set to empty string if not found
	}
	if description, ok := jsonData["Description"]; ok {
		task.Description = strings.ToLower(string(fmt.Sprintf("%v", description)))
	} else {
		task.Description = "" // Set to empty string if not found
	}

	//task.ExpiredDatetime = jsonData["expiredDatetime"].(*time.Time)

	var expiredTimePtr *time.Time
	if expiredTime, ok := jsonData["expiredDatetime"]; ok {
		if expiredTime != nil {
			tempExpiredTime := expiredTime.(time.Time)
			expiredTimePtr = &tempExpiredTime
		}
	}
	task.ExpiredDatetime = expiredTimePtr

	if pointsChange, ok := jsonData["PointsChange"]; ok {
		pointsChangeStr := string(fmt.Sprintf("%v", pointsChange))
		if pointsChangeInt, err := strconv.Atoi(pointsChangeStr); err == nil {
			task.PointsChange = pointsChangeInt
		}
	}

	task.Status = 0

	DB := db.Manager()

	//insert statement will not log in db-log , db-log  only log select //bug
	sql := DB.ToSQL(func(tx *gorm.DB) *gorm.DB { return tx.Create(&task) })
	lib.Log("sql:"+sql, "-db-log")

	result := DB.Create(&task)

	if result.Error != nil {
		log.Println("insert error:", result.Error)
		return c.JSON(http.StatusInternalServerError, lib.PyDict{
			"success": false,
			"message": "Failed to insert task",
		})
	}

	returnJson := lib.PyDict{
		"success": true,
		"message": "Task added successfully",
	}

	return c.JSON(http.StatusOK, returnJson)
}

// this is function to return tasks in db
func Tasks(c echo.Context) error {

	pageStr := c.QueryParam("page")
	page := 1
	if pageStr != "" {
		var err error
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 0 {
			page = 1
		}
	}
	pageSize := 10 // Adjust page size as needed
	pageSizeStr := c.QueryParam("rows")
	if pageSizeStr != "" {
		var err error
		pageSize, err = strconv.Atoi(pageSizeStr)
		if err != nil {
			pageSize = 10
		}
	}
	//set default value to select
	status := -10
	statusStr := strings.TrimSpace(c.QueryParam("status"))
	lib.Log("status:" + statusStr)
	if statusStr != "" {
		var err error
		status, err = strconv.Atoi(statusStr)
		if err != nil {
			status = -10
		}
	}

	offset := (page - 1) * pageSize
	log.Println("status:", status)
	lib.Log("status:" + fmt.Sprintf("%v", status))

	taskIDoruserID := strings.TrimSpace(strings.ToLower(c.Param("taskIDoruserID")))

	taskList := []db.Task{}

	DB := db.Manager()

	query := DB

	if strings.TrimSpace(taskIDoruserID) != "" {
		query = query.Where("task_id = ? OR user_id = ?", taskIDoruserID, taskIDoruserID)
		//query = query.Where("user_id = ?", strings.TrimSpace(taskIDoruserID))
	}

	if status != -10 {
		query = query.Where("status = ?", status)
	}

	// Get the total count of query results first
	var total int64
	tempQuery := query
	tempQuery = tempQuery.Model(&db.Task{})
	tempQuery.Count(&total)
	sql := tempQuery.ToSQL(func(tx *gorm.DB) *gorm.DB { return tx.Count(&total) })
	lib.Log("sql:"+sql, "-db-log")

	query = query.Limit(pageSize).Offset(offset)
	log.Println("sql:", query.Statement.SQL.String())
	sql = query.ToSQL(func(tx *gorm.DB) *gorm.DB { return tx.Find(&taskList) })
	lib.Log("sql:"+sql, "-db-log")
	//lib.Log("sql:", tempQuery.Statement.SQL.String())

	result := query.Find(&taskList)
	if result.Error != nil {
		log.Println("select error:", result.Error)
		return c.JSON(http.StatusInternalServerError, lib.PyDict{
			"success": false,
			"message": "Failed to retrieve tasks",
		})
	}

	returnJson := lib.PyDict{
		"success":       true,
		"message":       "Tasks retrieved successfully",
		"data":          taskList,
		"page":          page,
		"rows_per_page": pageSize,
		"total":         total,
	}
	//log the response returnJson jsonencode
	returnJsonEncoded, err := json.Marshal(returnJson)
	if err != nil {
		log.Println("json encode error:", err)
		return c.JSON(http.StatusInternalServerError, lib.PyDict{
			"success": false,
			"message": "Failed to encode JSON response",
		})
	}
	lib.Log(string(returnJsonEncoded))
	//lib.Log()
	return c.JSON(http.StatusOK, returnJson)
}

func RunTask(rows ...int) error {
	log.Println("runTask")
	//now := time.Now().UTC()
	// := now.Format("2006-01-02 15:04:05")
	//formatted := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	//lib.Log("======runTask" + formatted + "======")
	rowNum := 1
	if len(rows) > 0 {
		if rows[0] > 0 {
			rowNum = rows[0]
		}
	}

	workId := 1

	workIdStr := os.Getenv("workId")
	var err error
	if workIdStr == "" {
		workIdStr = "1" // Default port if PORT is not set
		fmt.Println("Environment variable workId not set, using default workId 1")
	} else {
		workId, err = strconv.Atoi(workIdStr)
		if err != nil {
			workId = 1
			fmt.Println("Environment variable workId is not a valid integer, using default workId 1")
		}
	}

	if len(rows) > 1 {
		if rows[1] > 0 {
			workId = rows[1]
		}
	}

	//lib.Log("======workId:" + fmt.Sprintf("%v", workId) + "======")

	offset := (workId - 1) * rowNum
	//lib.Log("======offset:" + fmt.Sprintf("%v", offset) + "======")
	//query the task limit by rowNum , and order by id
	DB := db.Manager()

	taskList := []db.Task{}
	result := DB.Where("Status = 0").Order("id").Limit(rowNum).Offset(offset).Find(&taskList)
	//use mutex lock to avoid duplicate select task
	var mutex = &sync.Mutex{}
	mutex.Lock()
	if result.Error != nil {
		log.Println("select error:", result.Error)
		lib.Log("select error:" + fmt.Sprintf("%v", result.Error))
		mutex.Unlock()
		return result.Error
	}

	// Start a new transaction
	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		log.Println("Transaction begin error:", err)
		lib.Log("select error:" + fmt.Sprintf("%v", err))
		mutex.Unlock()
		return err
	}

	for _, task := range taskList {
		uID := task.UserID
		lib.Log("======uID:" + fmt.Sprintf("%v", uID) + "======")
		var userPointBalance db.UserPointBalance
		// 1. SELECT * FROM userPointBalance WHERE user_id = uID FOR UPDATE
		txResult := tx.Set("gorm:query_option", "FOR UPDATE").Where("user_id = ?", uID).First(&userPointBalance)
		if txResult.Error != nil {
			if txResult.Error == gorm.ErrRecordNotFound {
				lib.Log("userPointBalance not find  add record:" + fmt.Sprintf("%v", uID))
				userPointBalance = db.UserPointBalance{
					UserID:  uID,
					Balance: 0,
				}
				txResult = tx.Create(&userPointBalance)
			}
			if txResult.Error != nil {
				lib.Log("SELECT user point error or create(LINE 289):" + fmt.Sprintf("%v", txResult.Error))
				tx.Rollback()
				//update task status -1
				TaskFailedTransaction(int(task.ID), 293)
				break //break the loop
			}
		}
		lib.Log("userPointBalance:" + fmt.Sprintf("%v", userPointBalance))

		// 2. UPDATE userPointBalance SET balance = balance + task.PointsChange WHERE user_id = uID
		newBalance := userPointBalance.Balance + int64(task.PointsChange)
		if txResult := tx.Model(&db.UserPointBalance{}).Where("user_id = ?", uID).Update("balance", newBalance); txResult.Error != nil {
			lib.Log("UPDATE user point error(LINE 305):" + fmt.Sprintf("%v", txResult.Error))
			tx.Rollback()
			//update task status -1
			TaskFailedTransaction(int(task.ID), 304)
			break //break the loop
		}

		if newBalance < 0 {
			// 3. if balance < 0 rollback and update task status=1
			tx.Rollback()
			lib.Log("userPointBalance <0 rollback and update task status=1 (LINE 318):" + fmt.Sprintf("%v", newBalance))
			//update task status -1
			TaskFailedTransaction(int(task.ID), 313)
			break //break the loop
		}

		if task.PointsChange < 0 {
			//4. Search the transaction_record  where user_id = uID order by ExpiredDatetime  DESC, and update the UsedPoints each records if records' PointsChange 10  then task.PointsChange -100 ,used should be -10 , then fetch next until used 100
			pointsToDeduct := task.PointsChange
			var transactionRecords []db.TransactionRecord

			forBreak := true
			//while loop
			for {
				// Infinite loop
				//points_change > 0 means plus  and (used_points + points_change) > 0 means has balnace to use
				if txResult := tx.Where("user_id = ?", uID).Where("points_change >= ?", 0).Where("(used_points + points_change) > 0").Order("expired_datetime DESC").Limit(10).Find(&transactionRecords); txResult.Error != nil {
					lib.Log("search transaction_record error(LINE 338):" + fmt.Sprintf("%v", txResult.Error))
					tx.Rollback()
					//update task status -1
					TaskFailedTransaction(int(task.ID), 331)
					forBreak = false
					break
				}

				doUpdatedRecord := true
				for _, record := range transactionRecords {
					if pointsToDeduct >= 0 {
						//break records for
						break
					}
					//this should not happen
					if record.PointsChange <= 0 {
						continue
					}
					//pointsToDeduct is minus and record.PointsChange is plus
					availablePoints := pointsToDeduct + record.PointsChange + record.UsedPoints
					deductAmount := 0
					//deductAmount is minus
					//availablePoints >= 0 , means record.PointsChange used pointsToDeduct
					if availablePoints >= 0 {
						deductAmount = pointsToDeduct
					}

					//if availablePoints < 0  , means record.PointsChange all used, -(record.PointsChange + record.UsedPoints)
					if availablePoints < 0 {
						deductAmount = -(record.PointsChange + record.UsedPoints)
					}
					lib.Log("update transaction_record used_points (LINE 359):" + fmt.Sprintf("%v", record))
					if txResult := tx.Model(&db.TransactionRecord{}).Where("id = ?", record.ID).Updates(map[string]interface{}{
						"used_points": int(record.UsedPoints) + deductAmount,
						"description": gorm.Expr("CONCAT(description, ?)", "\nUsed:"+task.TaskID),
					}); txResult.Error != nil {
						lib.Log("update transaction_record used_points error(LINE 364):" + fmt.Sprintf("%v", txResult.Error))
						tx.Rollback()
						//update task status -1
						TaskFailedTransaction(int(task.ID), 374)
						forBreak = false
						doUpdatedRecord = false
						//break records for
						break
					}

					//balance to deduct
					pointsToDeduct = pointsToDeduct - deductAmount
					lib.Log(fmt.Sprintf("deduct from record %d: used points: %d points, new pointsToDeduct: %d", record.ID, deductAmount, pointsToDeduct))
				}

				if doUpdatedRecord != true {
					lib.Log("update transaction_record used_points error !!!(LINE 392)")
					forBreak = false
					break
				}

				if pointsToDeduct >= 0 {
					lib.Log("Finshed deduct from transaction records!(LINE 399)")
					// Optionally, you can rollback or take other actions here
					//break for
					break
				}
			}

			if forBreak != true {
				lib.Log("update transaction_record used_points error !!!(LINE 407)")
				break
			}
		}

		if txResult := tx.Model(&db.Task{}).Where("id = ?", task.ID).Update("status", 1); txResult.Error != nil {
			lib.Log("update task status = 1  error(LINE 416):" + fmt.Sprintf("%v", txResult.Error))
			tx.Rollback()
			//update task status -1
			TaskFailedTransaction(int(task.ID), 402)
			break //break the loop
		}

		// 	UserID          string     `gorm:"column:user_id;type:char(64);not null"`
		// PointsChange    int        `gorm:"not null"`
		// TransactionDate time.Time  `gorm:"column:transaction_date;type:datetime;default:CURRENT_TIMESTAMP"`
		// Description     string     `gorm:"type:varchar(255)"`
		// TransactionType string     `gorm:"type:varchar(255)"`
		// TaskID          string     `gorm:"column:task_id;type:char(64)"`
		// ExpiredDatetime *time.Time `gorm:"column:expired_datetime;type:datetime"`
		// 5. INSERT INTO transaction_records (user_id, points_change, transaction_date, description, task_id)
		transactionRecord := db.TransactionRecord{
			UserID:          uID,
			PointsChange:    task.PointsChange,
			TransactionDate: time.Now(),
			Description:     task.Description,
			TaskID:          task.TaskID,
			ExpiredDatetime: task.ExpiredDatetime,
			TransactionType: task.TaskType,
		}
		if txResult := tx.Create(&transactionRecord); txResult.Error != nil {
			lib.Log("insert transaction_records error(LINE 440):" + fmt.Sprintf("%v", txResult.Error))
			tx.Rollback()
			//update task status -1
			TaskFailedTransaction(int(task.ID), 427)
			break //break the loop
		}

		lib.Log("======runTask" + fmt.Sprintf("%v", newBalance) + "======")
		// 5. Commit the transaction
		if err := tx.Commit().Error; err != nil {
			lib.Log("commit error(LINE 452)" + fmt.Sprintf("%v", err))
			tx.Rollback()
			break
		}
		//here maybe do call back
		//webhook_url //https://9000-idx-helloworld-1725508396294.cluster-qpa6grkipzc64wfjrbr3hsdma2.cloudworkstations.dev/?monospaceUid=363848
		// fmt.Println("Example 2: POST request with JSON body and headers")
		// jsonBody := []byte(`{"title": "foo", "body": "bar", "userId": 1}`)
		// headers := map[string]string{"Content-Type": "application/json"}
		// postResp, postErr := lib.HttpRequest("https://jsonplaceholder.typicode.com/posts", http.MethodPost, headers, bytes.NewBuffer(jsonBody))
		// if postErr != nil {
		// 	log.Fatalf("POST request failed: %v", postErr)
		// }
		// fmt.Printf("POST Response:\n%s\n\n", postResp)
		postDict := lib.PyDict{
			"success": true,
			"message": "Task added successfully",
			"data":    task,
		}
		//test url
		webhook_url := "http://localhost:9002/api/print_all"
		headers := map[string]string{"Content-Type": "application/json"}
		returnJsonEncoded, _ := json.Marshal(postDict)
		jsonBody := []byte(string(returnJsonEncoded))
		lib.Log("url:" + string(webhook_url))
		postResp, postErr := lib.HttpRequest(webhook_url, http.MethodPost, headers, bytes.NewBuffer(jsonBody))
		if postErr != nil {
			lib.Log("POST request failed:" + fmt.Sprintf("%v", postErr))
		}
		//fmt.Printf("POST Response:\n%s\n\n", postResp)
		lib.Log("POST Response v  postResp:" + fmt.Sprintf("%v", postResp))
		lib.Log("POST Response postResp:" + fmt.Sprintf("POST Response:\n%s\n", postResp))
		postRespString := string(postResp)
		lib.Log(fmt.Sprintf("POST Response postRespString:\n%s\n", postRespString))
	}
	//lib.Log("======runTask" + fmt.Sprintf("%v", taskList) + "======")
	mutex.Unlock()
	//lib.Log("======runTask" + fmt.Sprintf("%v", rowNum) + "======")
	return nil
}

func TaskFailedTransaction(taskID int, lineNume int) error {
	DB := db.Manager()
	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		log.Println("Transaction begin error:", err)
		lib.Log("select error:" + fmt.Sprintf("%v", err))
		return err
	}
	//select task for update
	txTaskResult := tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ?", taskID).First(&db.Task{})
	if txTaskResult.Error != nil {
		lib.Log("Task not find  add record:" + fmt.Sprintf("%v", taskID))
		lib.Log("Error:" + fmt.Sprintf("%v", txTaskResult.Error))
	} else {
		//update task status -1
		if txTaskResult = tx.Model(&db.Task{}).Where("id = ?", taskID).Update("status", -1); txTaskResult.Error != nil {
			lib.Log("update task status = -1  error(LINE " + fmt.Sprintf("%v", lineNume) + " ):" + fmt.Sprintf("%v", txTaskResult.Error))
			tx.Rollback()
		}
	}
	return nil
}

// collect the expired points_change > 0 records and put into Task to minus
func ExpiredTask() error {

	records := []db.TransactionRecord{}
	DB := db.Manager()
	now := time.Now()
	// expired_datetime is not null  ,points_change > 0 has expired_datetime // (used_points + points_change) > 0 means balance > 0
	result := DB.Where("expired_datetime is not null").Where("expired_datetime < ?", now).Where("(used_points + points_change) > 0").Find(&records)
	//result := DB.Where("expired_datetime is not null").Where("expired_datetime <") //  and points_change > 0 and (used_points + points_change) > 0 ", now).Find(&records)
	//log

	if result.Error != nil {
		log.Println("select error:", result.Error)
		lib.Log("select error:" + fmt.Sprintf("%v", result.Error))
		return result.Error
	}

	//fetch query
	for _, record := range records {
		task := db.Task{
			TaskID:          record.TaskID,
			UserID:          record.UserID,
			TaskType:        "Expired",
			Description:     record.Description + " Expired ",
			PointsChange:    -record.PointsChange,
			ExpiredDatetime: nil,
		}
		result := DB.Create(&task)
		if result.Error != nil {
			log.Println("insert error:", result.Error)
			lib.Log("insert error:" + fmt.Sprintf("%v", result.Error))
			return result.Error
		}
	}

	return nil
}

func GetUserPointBalancesWithTransactions() error {
	DB := db.Manager()

	// Define the structs to hold the joined data
	type Result struct {
		db.UserPointBalance
		TransactionRecords []db.TransactionRecord
	}

	// Execute the query with JOIN
	var results []Result
	err := DB.Model(&db.UserPointBalance{}).
		Select("user_point_balances.*, transaction_records.*").
		Joins("LEFT JOIN transaction_records ON transaction_records.user_id = user_point_balances.user_id").
		Scan(&results).Error

	if err != nil {
		log.Println("select join error:", err)
		lib.Log("select join error:" + fmt.Sprintf("%v", err))
		return err
	}

	// Process the results
	for _, result := range results {
		// Log the user point balance
		lib.Log("User Point Balance:")
		lib.Log(fmt.Sprintf("  ID: %d", result.UserPointBalance.ID))
		lib.Log(fmt.Sprintf("  UserID: %s", result.UserPointBalance.UserID))
		lib.Log(fmt.Sprintf("  Balance: %d", result.UserPointBalance.Balance))
		//lib.Log(fmt.Sprintf("  CreateDatetime: %s", result.UserPointBalance.CreateDatetime.String()))
		//lib.Log(fmt.Sprintf("  UpdateDatetime: %s", result.UserPointBalance.UpdateDatetime.String()))

		// Log associated transaction records
		lib.Log("  Transaction Records:")
		for _, transaction := range result.TransactionRecords {
			lib.Log("    Transaction Record:")
			lib.Log(fmt.Sprintf("      ID: %d", transaction.ID))
			lib.Log(fmt.Sprintf("      UserID: %s", transaction.UserID))
			lib.Log(fmt.Sprintf("      PointsChange: %d", transaction.PointsChange))
			lib.Log(fmt.Sprintf("      TransactionDate: %s", transaction.TransactionDate.String()))
			lib.Log(fmt.Sprintf("      Description: %s", transaction.Description))
			lib.Log(fmt.Sprintf("      TransactionType: %s", transaction.TransactionType))
			lib.Log(fmt.Sprintf("      TaskID: %s", transaction.TaskID))
			lib.Log(fmt.Sprintf("      ExpiredDatetime: %s", transaction.ExpiredDatetime.String()))
			//lib.Log(fmt.Sprintf("      CreateDatetime: %s", transaction.CreateDatetime.String()))
			//lib.Log(fmt.Sprintf("      UpdateDatetime: %s", transaction.UpdateDatetime.String()))
		}
	}
	return nil
}
