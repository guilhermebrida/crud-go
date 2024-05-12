package main

import (
	query "awesomeProject/db"
	"fmt"
	//_ "github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "front/templates/index.html")
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	db, err := query.ConnectDB()

	res, err := query.SelectCadastro(db)
	for res.Next() {
		var nome string
		if err := res.Scan(&nome); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Nome: %s\n", nome)
	}
	if err != nil {
		return
	}

	http.HandleFunc("/", index)

	http.HandleFunc("/users", query.SelectObject(db))

	http.HandleFunc("/names", query.SelectNames(db))

	http.HandleFunc("/insert", query.InsertPerson(db))

	http.HandleFunc("/add", query.Add())

	log.Fatal(http.ListenAndServe(":8080", nil))

}
