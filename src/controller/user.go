package controller

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"web/src/db"
	"web/src/lib"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Users(c echo.Context) error {

	userIdStr := c.Param("userID")
	userID := ""
	if strings.TrimSpace(userIdStr) != "" {
		userID = strings.TrimSpace(strings.ToLower(userIdStr))
	}

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
	offset := (page - 1) * pageSize
	log.Println("offset:", offset)

	userList := []db.UserPointBalance{}

	DB := db.Manager()

	query := DB
	if strings.TrimSpace(userID) != "" {
		query = query.Where("user_id = ?", userID)
	}
	// Get the total count of query results first
	var total int64
	tempQuery := query
	tempQuery = tempQuery.Model(&db.UserPointBalance{})
	tempQuery.Count(&total)
	sql := tempQuery.ToSQL(func(tx *gorm.DB) *gorm.DB { return tx.Count(&total) })
	lib.Log("sql:", tempQuery.Statement.SQL.String())

	query = query.Limit(pageSize).Offset(offset)
	result := query.Find(&userList)
	sql = query.ToSQL(func(tx *gorm.DB) *gorm.DB { return tx.Find(&userList) })
	lib.Log("sql:"+sql, "-db-log")

	if result.Error != nil {
		log.Println("select error:", result.Error)
		return c.JSON(http.StatusInternalServerError, lib.PyDict{
			"success": false,
			"message": "Failed to retrieve users",
		})
	}

	returnJson := lib.PyDict{
		"success":       true,
		"message":       "Users retrieved successfully",
		"data":          userList,
		"page":          page,
		"rows_per_page": pageSize,
		"total":         total,
	}

	return c.JSON(http.StatusOK, returnJson)

}

func Records(c echo.Context) error {

	userIdStr := c.Param("userID")
	userID := ""
	if strings.TrimSpace(userIdStr) != "" {
		userID = strings.TrimSpace(strings.ToLower(userIdStr))
	}

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

	offset := (page - 1) * pageSize
	log.Println("offset:", offset)

	recordList := []db.TransactionRecord{}

	DB := db.Manager()

	query := DB
	if strings.TrimSpace(userID) != "" {
		query = query.Where("user_id = ?", userID)
	}

	// Get the total count of query results first
	var total int64
	tempQuery := query
	tempQuery = tempQuery.Model(&db.TransactionRecord{})
	tempQuery.Count(&total)
	sql := tempQuery.ToSQL(func(tx *gorm.DB) *gorm.DB { return tx.Count(&total) })
	lib.Log("sql:", tempQuery.Statement.SQL.String())

	query = query.Limit(pageSize).Offset(offset)
	//fetch query
	result := query.Find(&recordList)
	sql = query.ToSQL(func(tx *gorm.DB) *gorm.DB { return tx.Find(&recordList) })
	lib.Log("sql:"+sql, "-db-log")

	if result.Error != nil {
		log.Println("select error:", result.Error)
		return c.JSON(http.StatusInternalServerError, lib.PyDict{
			"success": false,
			"message": "Failed to retrieve records",
		})
	}

	returnJson := lib.PyDict{
		"success":       true,
		"message":       "Records retrieved successfully",
		"data":          recordList,
		"page":          page,
		"rows_per_page": pageSize,
		"total":         total,
	}

	return c.JSON(http.StatusOK, returnJson)
}

func PrintPayload(c echo.Context) error {
	//check post data type is json
	if c.Request().Header.Get("Content-Type") == "application/json" {
		var json map[string]interface{} = map[string]interface{}{}
		if err := c.Bind(&json); err != nil {
			return c.JSON(http.StatusBadRequest, lib.PyDict{
				"success": false,
				"message": "Invalid JSON payload",
			})
		}
		return c.JSON(http.StatusOK, json)
	}
	bodyContent, _ := io.ReadAll(c.Request().Body)
	return c.String(http.StatusOK, fmt.Sprintf("%v", string(bodyContent)))
}
