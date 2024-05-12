package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	_ "html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type Person struct {
	ID            int    `json:"id"`
	Nome          string `json:"name"`
	Sobrenome     string `json:"last-name"`
	CPF           string `json:"cpf"`
	ReceptionDate string `json:"reception_date"`
	Email         string `json:"email"`
}

func ConnectDB() (*sql.DB, error) {
	connStr := "user=postgres dbname=postgres password=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("Conexão bem-sucedida com o banco de dados")
	return db, nil
}

func SelectCadastro(db *sql.DB) (*sql.Rows, error) {
	// Executar uma consulta SQL
	rows, err := db.Query("SELECT \"Nome\" FROM cadastro")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var nome string
		if err := rows.Scan(&nome); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Nome: %s\n", nome)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return rows, nil
}

func SelectObject(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM cadastro")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var people []Person
		for rows.Next() {
			var person Person
			if err := rows.Scan(&person.Nome, &person.Sobrenome, &person.Email, &person.CPF, &person.ReceptionDate, &person.ID); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			people = append(people, person)
		}

		//sendJson(w, people)
		serveHTML(w, people, "templates/users.html")

	}
}

func SelectNames(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT \"Nome\" FROM cadastro")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var names []string
		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			names = append(names, name)
		}

		//sendJson(w, names)
		serveHTML(w, names, "templates/names.html")

		//Parse do template HTML
		//tmpl, err := template.ParseFiles("templates/names.html")
		//if err != nil {
		//	http.Error(w, err.Error(), http.StatusInternalServerError)
		//	return
		//}
		//
		//// Renderize o template com os dados e envie a resposta
		//if err := tmpl.Execute(w, names); err != nil {
		//	http.Error(w, err.Error(), http.StatusInternalServerError)
		//	return
		//}
		//http.ServeFile(w, r, "templates/names.html")
	}
}

func Add() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serveHTML(w, "", "templates/index.html")
		fmt.Printf(r.FormValue("form"))
	}
}

func InsertPerson(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Requisição recebida para inserção")
		fmt.Printf(r.Host)
		fmt.Printf(r.Proto)
		// Lê e imprime o corpo da requisição
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Erro ao ler o corpo da requisição", http.StatusInternalServerError)
			log.Println("Erro ao ler o corpo da requisição:", err)
			return
		}
		defer r.Body.Close()

		fmt.Println("Corpo da Requisição:", string(body))
		var person Person
		if err := json.Unmarshal(body, &person); err != nil {
			http.Error(w, "Erro ao decodificar o corpo da requisição", http.StatusBadRequest)
			log.Println("Erro ao decodificar o corpo da requisição:", err)
			return
		}

		fmt.Println("Nome:", person.Nome)
		fmt.Println("Email:", person.Email)
		fmt.Println("Last Name:", person.Sobrenome)
		fmt.Println("CPF:", person.CPF)

		// Preparando a declaração SQL de inserção
		stmt, err := db.Prepare("INSERT INTO cadastro (\"Nome\", \"Sobrenome\", \"Email\", \"CPF\") VALUES ($1, $2, $3, $4)")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println("Erro ao preparar declaração SQL:", err)
			return
		}
		defer stmt.Close()

		// Executando a instrução de inserção com os dados recebidos
		_, err = stmt.Exec(person.Nome, person.Sobrenome, person.Email, person.CPF)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println("Erro ao executar instrução de inserção:", err)
			return
		}

		log.Println("Inserção realizada com sucesso")

		// Retornando uma resposta de sucesso
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Inserção realizada com sucesso!"))
	}
}

func sendJson(w http.ResponseWriter, data interface{}) {
	// Codificando os resultados em JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func serveHTML(w http.ResponseWriter, data interface{}, fileName string) {
	tmpl, err := template.ParseFiles(fileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Renderize o template com os dados e envie a resposta
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
