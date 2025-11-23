package main

import (
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 题目一
type students struct {
	gorm.Model
	Name  string
	Age   int
	Grade string
}

func getConnect() *gorm.DB {
	dsn := "root:root@tcp(127.0.0.1:3306)/go_test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func createTable() {
	db := getConnect()
	var students students
	db.AutoMigrate(&students)
}

func insertData(s *students) {
	db := getConnect()
	db.Create(s)
}

func selectData(s *students) {
	db := getConnect()
	db.Debug().Where("age > 18").Find(&s)
}

func updateData(s *students) {
	db := getConnect()
	db.Debug().Model(s).Where("name = ?", "张三").Update("grade", "四年级")
}

func deleteData(s *students) {
	db := getConnect()
	db.Debug().Where("age < ?", 15).Delete(s)
}

// 题目二
type account struct {
	gorm.Model
	Balance int
}

type transaction struct {
	gorm.Model
	FromAccountID uint
	ToAccountID   uint
	Amount        int
}

func transferAmountMethod(fromId uint, toId uint, amount int, tx *gorm.DB) error {
	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}
	accountA := &account{}
	accountB := &account{}
	//加锁查询账户信息
	tx.Debug().Clauses(clause.Locking{Strength: "UPDATE"}).Find(accountA, "id = ? ", fromId)
	tx.Debug().Clauses(clause.Locking{Strength: "UPDATE"}).Find(accountB, "id = ? ", toId)
	if accountA.Balance < amount {
		return errors.New("balance is not enough")
	}
	accountA.Balance -= amount
	tx.Debug().Model(accountA).Update("balance", accountA.Balance)
	accountB.Balance += amount
	tx.Debug().Model(accountB).Update("balance", accountB.Balance)
	//记录交易流水
	transaction := &transaction{
		FromAccountID: fromId,
		ToAccountID:   toId,
		Amount:        amount,
	}
	tx.Debug().Create(transaction)
	return nil
}

// sqlx 题目一
type employees struct {
	ID         uint           `db:"id"`
	CreatedAt  time.Time      `db:"created_at"`
	UpdatedAt  time.Time      `db:"updated_at"`
	DeletedAt  gorm.DeletedAt `db:"deleted_at"`
	Name       string         `db:"name"`
	Department string         `db:"department"`
	Salary     int            `db:"salary"`
}

func QueryEmployeesBySqlx(db *sqlx.DB, department string) ([]employees, error) {
	var emps []employees
	query := "SELECT * FROM employees WHERE department = ? and deleted_at IS NULL"
	err := db.Select(&emps, query, department)
	if err != nil {
		return nil, err
	}
	return emps, nil
}

func QueryEmployeesMaxSalaryBySqlx(db *sqlx.DB) (employees, error) {
	var emp employees
	query := "SELECT * FROM employees WHERE deleted_at IS NULL ORDER BY salary DESC LIMIT 1"
	err := db.Get(&emp, query)
	if err != nil {
		return emp, err
	}
	return emp, nil
}

type book struct {
	ID        uint           `db:"id"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
	DeletedAt gorm.DeletedAt `db:"deleted_at"`
	Title     string         `db:"title"`
	Author    string         `db:"author"`
	Price     float64        `db:"price"`
}

func QueryBooksByPriceRange(db *sqlx.DB, price float64) ([]book, error) {
	var books []book
	query := "SELECT * FROM books WHERE price > ? and deleted_at IS NULL"
	err := db.Select(&books, query, price)
	if err != nil {
		return nil, err
	}
	return books, nil
}

// 进阶gorm
// 题目一
type User struct {
	gorm.Model
	Name     string
	Email    string
	Password string
	Post_num int `gorm:"default:0"`
	Posts    []Post
}

type Post struct {
	gorm.Model
	Title    string
	Content  string
	UserID   uint
	State    string
	Comments []Comment
}

type Comment struct {
	gorm.Model
	Content string
	PostID  uint
}

// 题目二
func getOneUserInfo(db *gorm.DB, userID uint) (User, error) {
	var user User
	err := db.Preload("Posts").Preload("Posts.Comments").First(&user, userID).Error
	return user, err
}

// 题目三
func (p *Post) AfterCreate(tx *gorm.DB) (err error) {
	userId := p.UserID
	//同事查询并且更新用户的post_num字段+1 原子操作
	errNew := tx.Model(&User{}).Where("id", userId).UpdateColumn("post_num", gorm.Expr("post_num + ?", 1)).Error
	return errNew
}

func (c *Comment) AfterDelete(tx *gorm.DB) (err error) {
	postId := c.PostID
	var postCount int64
	errNew := tx.Debug().Model(&Comment{}).Where("post_id = ?", postId).Count(&postCount).Error
	if (postCount == 0) && (errNew == nil) {
		tx.Debug().Model(&Post{}).Where("id = ?", postId).Update("state", "无评论")
	}
	return errNew
}

func main() {
	//createTable() 创建表

	//题目一 {

	//插入数据
	// student := students{Name: "张三", Age: 20, Grade: "三年级"}
	// insertData(&student)

	//查询数
	// student := students{}
	// selectData(&student)
	// fmt.Println(student)
	// updateData(&student)

	//删除
	// deleteData(&student)
	//}

	//题目二 {

	// db := getConnect()
	// db.AutoMigrate(&book{})
	// db.AutoMigrate(&transaction{})

	// a := account{Balance: 100}
	// b := account{Balance: 500}
	// db.Create(&a)
	// db.Create(&b)

	// transferAmount := 100
	// db.Transaction(func(tx *gorm.DB) error {
	// 	return transferAmountMethod(1, 2, transferAmount, tx)
	// })
	// fmt.Print()

	//}

	//sqlx
	// db, err := sqlx.Connect("mysql", "root:root@tcp(127.0.0.1:3306)/go_test?charset=utf8mb4&parseTime=True&loc=Local")
	// if err != nil {
	// 	panic(err)
	// }
	//题目一
	// emps, err := QueryEmployeesBySqlx(db, "技术部")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(emps)

	// maxEmp, err := QueryEmployeesMaxSalaryBySqlx(db)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(maxEmp)

	//题目二
	// books, err := QueryBooksByPriceRange(db, 50.0)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(books)

	//进阶gorm
	//题目一
	db := getConnect()
	// db.AutoMigrate(&User{})
	// db.AutoMigrate(&Post{})
	// db.AutoMigrate(&Comment{})

	//题目二
	// user, err := getOneUserInfo(db, 1)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(user.Posts)

	//题目三
	// post := Post{
	// 	Title:   "第一篇文章",
	// 	Content: "这是我的第一篇文章内容",
	// 	UserID:  2,
	// }
	// db.Create(&post)

	//用的是 db.Debug().Delete(&Comment{}, 4) 这种按条件/主键删除的调用，
	// GORM 并不会先把整条记录加载到结构体再删除，因此传入到 AfterDelete
	// 的 c 结构体并没有被填充出 PostID 等字段，只有 Hook 收到的模型是零值，所以 c.PostID == 0
	// db.Debug().Delete(&Comment{}, 4)

	//所以要先查询再删除
	var comment Comment
	db.Debug().First(&comment, 5)
	db.Debug().Delete(&comment, 5)
}
