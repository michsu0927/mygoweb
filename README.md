This is a web server that announces whether or not a particular Go version has been tagged.

## Environment Variables

*   `DB_ENV`: Specifies the database type. Defaults to `sqlite`. If set to `mysql`, `DATABASE_DSN` must be configured. Valid values are `sqlite` or `mysql`.
*   `DATABASE_DSN`: Specifies the database connection string.
    *   Defaults to `data_db/test.db` (sqlite) if not set.
    *   If `DB_ENV` is set to `mysql`, this variable is **required**.
    *   Example (MySQL): `user:password@tcp(127.0.0.1:3306)/database_name?charset=utf8mb4&parseTime=True&loc=Local`
        *   `user`: The username for your MySQL database. 
        *   `password`: The password for your MySQL user.
        *   `tcp(127.0.0.1:3306)`: The network address and port of your MySQL server. `127.0.0.1` is the loopback address (localhost), and `3306` is the default MySQL port.
        *   `database_name`: The name of the database you want to connect to.
        *   `charset=utf8mb4`: Specifies the character set to use for communication with the MySQL server.
        *   `parseTime=True`: Tells the MySQL driver to parse `DATETIME` and `TIMESTAMP` columns from MySQL into `time.Time` objects in Go.
        *   `loc=Asia%2FTaipei`: Specifies the time zone to use for parsing time values , you can use `loc=Local` for the system's local time zone or use  `loc=Asia%2FTaipei` for the Taipei time zone.
    
    
*   `PORT`: Specifies the port for the server to listen on. Defaults to `3000` if not set.
*   `expireTime`: Specifies the time in `hh:mm` format when the `ExpiredTask` should run each day. Defaults to `00:00` if not set.expireTime  if not as `00:00`  it will not run.
*   `runTaskNum`: Specifies how many tasks `RunTask` should process at a time. Defaults to `1` if not set.
*   `runTaskSecs`: Specifies the interval in seconds between `RunTask` executions. Defaults to `1` second if not set.
*   `workId`: Specifies the `workId`. Defaults to `1` if not set.

    
    
# 中文

這是一個Go lang 的網路伺服器。

## 環境變數

*   `DB_ENV`: 指定資料庫類型。預設為 `sqlite`。如果設定為 `mysql`，則必須設定 `DATABASE_DSN`。有效值為 `sqlite` 或 `mysql`。
*   `DATABASE_DSN`: 指定資料庫連線字串。
    *   如果未設定，則預設為 `data_db/test.db` (sqlite)。
    *   如果 `DB_ENV` 設定為 `mysql`，則此變數為**必要**。
    *   範例 (MySQL): `user:password@tcp(127.0.0.1:3306)/database_name?charset=utf8mb4&parseTime=True&loc=Local`
        *   `user`: 您的 MySQL 資料庫使用者名稱。
        *   `password`: 您的 MySQL 使用者密碼。
        *   `tcp(127.0.0.1:3306)`: 您的 MySQL 伺服器的網路位址和連接埠。`127.0.0.1` 是迴路位址 (localhost)，`3306` 是預設的 MySQL 連接埠。
        *   `database_name`: 您要連線的資料庫名稱。
        *   `charset=utf8mb4`: 指定與 MySQL 伺服器通訊時使用的字元集。
        *   `parseTime=True`: 告訴 MySQL 驅動程式將 MySQL 中的 `DATETIME` 和 `TIMESTAMP` 資料行解析為 Go 中的 `time.Time` 物件。
        * `loc=Asia%2FTaipei`: 指定用於解析時間值的時區，可以使用`loc=Local` 系統的本地時區，也可以使用 `loc=Asia%2FTaipei` 台北時區。

*   `PORT`: 指定伺服器監聽的連接埠。如果未設定，則預設為 `3000`。
*   `expireTime`: 指定 `ExpiredTask` 每天應執行的時間，格式為 `hh:mm`。如果未設定，則預設為 `00:00`。如果不是 `00:00` 就不會執行 `ExpiredTask`。
*   `runTaskNum`: 指定 `RunTask` 一次應處理多少個任務。如果未設定，則預設為 `1`。
*   `runTaskSecs`: 指定 `RunTask` 執行之間的間隔秒數。如果未設定，則預設為 `1` 秒。
*   `workId`: 指定 `workId`。如果未設定，則預設為 `1`。