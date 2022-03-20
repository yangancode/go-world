package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"sync"
	"time"
)

type userBalance struct {
	Uid     int
	Name    string
	Balance int
	Version int
}

func (userBalance) TableName() string {
	return "user_balance"
}

// mysql update 是行锁，但并发会慢SQL
func initDB(dsn string) (*gorm.DB, error) {
	slowLogger := logger.New(
		//将标准输出作为Writer
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			//设定慢查询时间阈值为1ms
			SlowThreshold: 1 * time.Second,
			//设置日志级别，只有Warn和Info级别会输出慢查询日志
			LogLevel: logger.Warn,
		},
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: slowLogger})
	if err != nil {
		return nil, err
	}
	return db, err
}

func optimisticLockTest(db *gorm.DB) error {
	balance := userBalance{}
	tx := db.Where("uid=?", 1).Find(&balance)
	if err := tx.Error; err != nil {
		return err
	}
	fmt.Println("balance", balance)
	return nil
}

// gorm DB
func optimisticLockByVersion(db *gorm.DB, idx int) error {
	tx := db.Exec("update `user_balance` set `balance`=20, `version`=`version`+1 where uid = 1 and version=0")
	if err := tx.Error; err != nil {
		return err
	}
	if tx.RowsAffected == 1 {
		fmt.Printf("affect: %v, %v\n", idx, tx.RowsAffected)
	}
	return nil
}

func optimisticLockByCAS(db *gorm.DB, idx int) error {
	tx := db.Exec("update `user_balance` set `balance`=20 where uid = 1 and balance=100")
	if err := tx.Error; err != nil {
		return err
	}
	if tx.RowsAffected == 1 {
		fmt.Printf("affect: %v, %v\n", idx, tx.RowsAffected)
	}
	return nil
}

func main() {
	dsn := "root:123456@tcp(127.0.0.1:33306)/micro"
	db, err := initDB(dsn)
	if err != nil {
		panic(err)
	}
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		idx := i
		fmt.Println("idx", i)
		go func(idx int) {
			err = optimisticLockByCAS(db, idx)
			if err != nil {
				fmt.Println("err", err)
			}
			wg.Done()
		}(idx)
	}
	wg.Wait()
}
