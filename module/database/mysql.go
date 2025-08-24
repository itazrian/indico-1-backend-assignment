package database

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    _ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Init() {
    host := getKoneksi("MYSQL_HOST", "localhost") // local
    //host := getKoneksi("MYSQL_HOST", "mysql8") // docker
    
    port := getKoneksi("MYSQL_PORT", "3306")
    
    user := getKoneksi("MYSQL_USER", "root") // local
    pass := getKoneksi("MYSQL_PASSWORD", "") // local
    //user := getKoneksi("MYSQL_USER", "user") // docker
    //pass := getKoneksi("MYSQL_PASSWORD", "password") // docker 
    
    dbname := getKoneksi("MYSQL_DB", "test_indico_220825")

    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, port, dbname)

    var err error
    DB, err = sql.Open("mysql", dsn)
    if err != nil {
        log.Fatalf("Gagal Koneksi ke database: %v", err)
    }
    if err = DB.Ping(); err != nil {
        log.Fatalf("Gagal ping ke database: %v", err)
    }
    log.Println("Koneksi ke mysql berhasil")
}

func getKoneksi(k, def string) string {
    if v := os.Getenv(k); v != "" {
        return v
    }
    return def
}
