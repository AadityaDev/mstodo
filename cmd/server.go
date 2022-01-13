/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	v "github.com/aadityadev/mstodo/pkg/api/v1"
	"golang.org/x/xerrors"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	_ "github.com/go-sql-driver/mysql"
)

type server struct {
	v.UnimplementedUserServiceServer
}

type TodoItemModel struct{
	Id int `json:"Id"`
	Description string `json:"Description"`
	Completed bool `json:"Completed"`
}

func (*server) GetUser(ctx context.Context, request *v.UserRequest) (*v.UserResponse, error) {
	response := &v.UserResponse{};

	if request == nil {
		fmt.Println("request must not be nil")
		return response, xerrors.Errorf("request must not be nil")
	}
	id := request.Id
	fmt.Printf("User ID is: %d", id)

	name := request.Name;
	fmt.Println("User Name is: %d", name)

	response = &v.UserResponse{
		Name: "John Doe",
	}
	return response, nil
}

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("server called")

		if len(os.Args) != 2 {
			log.Fatal("Please specify to start the server or client.")
		}
	
		osArgs := os.Args[1:]
		fmt.Println(osArgs)
		
		address := "0.0.0.0:8000"
		
		lis, err := net.Listen("tcp", address)
		http.HandleFunc("/getTodo", GetTodo)
		http.HandleFunc("/createTodo", CreateTodo)
		http.HandleFunc("/updateTodo", UpdateTodo)
		http.HandleFunc("/", HelloServer)
		http.ListenAndServe(":8080", nil)
		
		db, err := sql.Open("mysql", "root:12345678@tcp(127.0.0.1:3306)/testdb")
		defer db.Close()

		if err != nil {
			log.Fatal(err)
		}

		sql := "SELECT * FROM testdb.todo;"
		res, err := db.Exec(sql)
	
		fmt.Printf("res %s", res)

		if err != nil {
			panic(err.Error())
		}

		if err != nil {
			log.Fatalf("Error %v", err)
		}

		fmt.Printf("Server is listening on %v ...\n", address)
		s := grpc.NewServer()
		v.RegisterUserServiceServer(s, &server{})
		s.Serve(lis)
	},
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var dat TodoItemModel
	da, err := ioutil.ReadAll(r.Body);
	if err != nil {
		log.Println(err)
		w.Header().Add("Content-Type", "application/json")
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.Unmarshal(da, &dat)
	fmt.Println("dat", dat);
	if err != nil {
		log.Println(err)
		w.Header().Add("Content-Type", "application/json")
		http.Error(w, err.Error(), 500)
		return
	}

	if dat.Description == "" {
		log.Println(err)
		w.Header().Add("Content-Type", "application/json")
		http.Error(w, "Description is missing", 500)
		return
	}

	db, err := sql.Open("mysql", "root:12345678@tcp(127.0.0.1:3306)/testdb")
	if err != nil {
		log.Println(err)
		w.Header().Add("Content-Type", "application/json")
		http.Error(w, err.Error(), 500)
		return
	}

	sql := "INSERT INTO testdb.todo(Description, Completed) VALUES (?,?);";
	res, err := db.Exec(sql, dat.Description, dat.Completed)
	if err != nil {
		log.Println(err)
		w.Header().Add("Content-Type", "application/json")
		http.Error(w, err.Error(), 500)
		return
	}

	newId, err := res.LastInsertId();
	defer db.Close()
	fmt.Println("res is ",newId)
	if err != nil {
		log.Println(err)
		w.Header().Add("Content-Type", "application/json")
		http.Error(w, err.Error(), 500)
		return
	}
	
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, "%d", newId)
}

func GetTodo(w http.ResponseWriter, r *http.Request) {
	dat, err := json.Marshal(r.Body)
	fmt.Println("dat %s", dat);
	fmt.Println("err %s", err);

	db, err := sql.Open("mysql", "root:12345678@tcp(127.0.0.1:3306)/testdb")
	if err != nil {
		log.Println(err)
		w.Header().Add("Content-Type", "application/json")
		http.Error(w, err.Error(), 500)
		return
	}

	sql := "SELECT * FROM testdb.todo;"
	res, err := db.Query(sql)
	defer res.Close()

	if err != nil {
		log.Println(err)
		w.Header().Add("Content-Type", "application/json")
		http.Error(w, err.Error(), 500)
		return
	}
	
	var todoItems []TodoItemModel;
	for res.Next() {
        var alb TodoItemModel
        if err := res.Scan(&alb.Id, &alb.Description, &alb.Completed); 
		err != nil {
            log.Println(err)
			w.Header().Add("Content-Type", "application/json")
			http.Error(w, err.Error(), 500)
			return
        }
        todoItems = append(todoItems, alb)
    }
    if err = res.Err(); 
	err != nil {
        log.Println(err)
		w.Header().Add("Content-Type", "application/json")
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Printf("%v\n", todoItems)

	b, err := json.Marshal(todoItems)
	if err != nil {
		log.Println(err)
		w.Header().Add("Content-Type", "application/json")
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
    fmt.Fprintf(w, "%s", string(b))
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	var res sql.Result
	var err error
	var dat TodoItemModel
	da, err := ioutil.ReadAll(r.Body);
	rowId, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64);
	fmt.Println("row id: ", rowId);

	if err != nil {
		log.Println(err)
		w.Header().Add("Content-Type", "application/json")
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.Unmarshal(da, &dat)
	fmt.Println("dat", dat);
	if err != nil {
		log.Println(err)
		w.Header().Add("Content-Type", "application/json")
		http.Error(w, err.Error(), 500)
		return
	}

	if rowId <= 0 {
		log.Println(err)
		w.Header().Add("Content-Type", "application/json")
		http.Error(w, "row id is missing!!!", 500)
		return
	}

	db, err := sql.Open("mysql", "root:12345678@tcp(127.0.0.1:3306)/testdb")
	if err != nil {
		log.Println(err)
		w.Header().Add("Content-Type", "application/json")
		http.Error(w, err.Error(), 500)
		return
	}


	if((dat.Description != "") && (dat.Completed == true || dat.Completed == false)) {
		sql := "UPDATE testdb.todo SET Description = ?, Completed = ? WHERE id = ?;";
		res, err = db.Exec(sql, dat.Description, dat.Completed, rowId)	
	} else if(dat.Description != "") {
		sql := "UPDATE testdb.todo SET Description = ? WHERE id = ?;";
		res, err = db.Exec(sql, dat.Description, rowId)
	} else if((dat.Completed == true || dat.Completed == false)) {
		sql := "UPDATE testdb.todo SET Completed = ? WHERE id = ?;";
		res, err = db.Exec(sql, dat.Completed, rowId)
	} 
	if err != nil {
		log.Println(err)
		w.Header().Add("Content-Type", "application/json")
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Println("result after update is: ", res);
	fmt.Println("result after update is: ", res);

	newId, err := res.RowsAffected();
	fmt.Println("result after affected is: ", newId);

	defer db.Close()
	fmt.Println("res is ",newId)
	if err != nil {
		log.Println(err)
		w.Header().Add("Content-Type", "application/json")
		http.Error(w, err.Error(), 500)
		return
	}
	
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, "%d", newId)
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
