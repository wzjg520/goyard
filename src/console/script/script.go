package main

//import (
//	//"database/sql"
//
//	//_ "github.com/go-sql-driver/mysql"
//	//"errors"
//	//"encoding/json"
//	//"database/sql"
//
//)
import (
	"db"
	"fmt"
	"encoding/json"
)

func main() {
	model := db.NewDb("mysql", "root:123456@tcp(localhost:3306)/jiajun")

	ret := model.Select("select * from test")

	defer model.Close()

	//fmt.Println(json.MarshalIndent(ret, "", "  "))
	//jsonByte, _ := json.Marshal(ret)
	jsonByte, _ := json.MarshalIndent(ret, "", "  ")

	fmt.Println(string(jsonByte))

	//db, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/jiajun")
	//checkError(err)
	//defer db.Close()
	//rows, err := db.Query("select * from test")
	//defer rows.Close()
	//checkError(err)
	//type personal struct {
	//	Id         int
	//	Name       string
	//	CreateTime int
	//	Email      string
	//}
	//
	////type personal map[string]interface{}
	//
	//rowsMap := make([]personal, 0)
	//
	//var id int
	//var name string
	//var createTime int
	//var email string
	//
	//for rows.Next() {
	//	rows.Scan(&id, &name, &createTime, &email)
	//	rowsMap = append(rowsMap, personal{
	//		Id:         id,
	//		Name:       name,
	//		CreateTime: createTime,
	//		Email:      email,
	//	})
	//	//rowsMap["id"] = id
	//	//rowsMap["name"] = name
	//	//rowsMap["createTime"] = createTime
	//	//rowsMap["email"] = email
	//	//rowsMap = append(rowsMap, personal{
	//	//	"id":id,
	//	//	"name":name,
	//	//	"createTime":createTime,
	//	//	"email":email,
	//	//})
	//}
	//jsonStr, err := json.MarshalIndent(rowsMap, "", "  ")
	//checkError(err)
	//fmt.Println(string(jsonStr))
}

//func checkError(err error) {
//	if err != nil {
//		log.Println(err)
//	}
//}
