package hello

import (
	"fmt"
	"net/http"

	"appengine"
	"database/sql"
	"encoding/json"

	"gopkg.in/gorp.v1"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	http.HandleFunc("/", handler)
}

type UserOne struct {
	Id       int64  `db:'Id' json:"id"`
	UserName string `db:'UserName'  json:"userName"`
	Email    string `db:'Email'  json:"email"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	db, err := sql.Open("mysql", "veeru:rahulhcs123*@tcp(173.194.239.224:3306)/testapp")
	
	checkErr(err, "DB Connection failed", c)

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{Engine: "InnoDB", Encoding: "UTF8"}}
	dbmap.AddTableWithName(UserOne{}, "users").SetKeys(true, "Id")
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Table creation failed", c)

	if err != nil {
		c.Errorf("Some error occurred %s", err)
		fmt.Fprint(w, "Some error occurred")
	} else {
		var users []UserOne
		_, err := dbmap.Select(&users, "select * from users")
		if err != nil {
			c.Errorf("Some error occurred %s", err)
			fmt.Fprint(w, "Some error occurred")
		} else {
			response, err := json.Marshal(users)
			checkErr(err, "Error while marshaling", c)
			fmt.Fprint(w, string(response))
		}
	}
}

func checkErr(err error, msg string, c appengine.Context) {
	if err != nil {
		c.Errorf("message :%v . The error is because %s", msg, err)
	}
}
