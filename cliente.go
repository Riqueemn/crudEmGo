package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Personagens struct {
	ID   int    `json:"id"`
	Nome string `json:"nome"`
}

//PersonagemHandler analina o request e delega para a função adequada
func PersonagemHandler(w http.ResponseWriter, r *http.Request) {
	sid := strings.TrimPrefix(r.URL.Path, "/personagens/") //Essa função retira toda a string anterior a /personagens/ incluindo /personagens/

	pid := make([]string, 2, 2)
	pid = strings.Split(sid, "=")

	ids := strings.Split(pid[0], ".")
	nome := string(pid[len(pid)-1])

	var id []int

	for i := 0; i < len(ids); i++ {
		a, _ := strconv.Atoi(ids[i])
		id = append(id, a)
	}

	switch {
	case r.Method == "GET" && (id[0] == 1 && id[1] > 0):
		fmt.Fprintf(w, "Visualizar")
		personagemPorID(w, r, id[1])
	case r.Method == "GET" && id[0] == 2:
		fmt.Fprintf(w, "Inserir")
		inserirPersonagem(w, r, id[1])
	case r.Method == "GET" && id[0] == 3:
		fmt.Fprintf(w, "Excluir")
		deletarPersonagem(w, r, id[1])
	case r.Method == "GET" && id[0] == 4:
		fmt.Fprintf(w, "Editar")
		editarPersonagem(w, r, nome, id[1])
	case r.Method == "GET":
		fmt.Fprintf(w, "Vizualizar todos")
		personagemTodos(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Desculpa, não encontramos sua requisição :(")
	}
}

func personagemPorID(w http.ResponseWriter, r *http.Request, id int) {
	db, err := sql.Open("mysql", "root:@/starwars")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var p Personagens
	db.QueryRow("SELECT id, nome from personagens where id = ?", id).Scan(&p.ID, &p.Nome) //Seleciona apenas uma única linha. Após, foi feito uma concatenação para mapear os valores para a variável p.

	json, _ := json.Marshal(p) //Gerando o json com a struct p

	w.Header().Set("Content-Type", "application/json") //Setando o Header
	fmt.Fprint(w, string(json))                        //Tranformando o json em string e passando para o responseWriter
}

func personagemTodos(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:@/starwars")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, _ := db.Query("SELECT id, nome from personagens")
	defer rows.Close()

	var p []Personagens
	for rows.Next() {
		var personagem Personagens
		rows.Scan(&personagem.ID, &personagem.Nome)
		p = append(p, personagem)
	}

	json, _ := json.Marshal(p)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(json))
}

func inserirPersonagem(w http.ResponseWriter, r *http.Request, id int) {
	db, err := sql.Open("mysql", "root:@/starwars")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, _ := db.Prepare("INSERT INTO personagens (id, nome) values(?, ?)")
	stmt.Exec(id, "Millenium Falcon")
	defer stmt.Close()

}

func deletarPersonagem(w http.ResponseWriter, r *http.Request, id int) {
	db, err := sql.Open("mysql", "root:@/starwars")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, _ := db.Prepare("DELETE from personagens where id=?")
	stmt.Exec(id)
	defer stmt.Close()
}

func editarPersonagem(w http.ResponseWriter, r *http.Request, nome string, id int) {
	db, err := sql.Open("mysql", "root:@/starwars")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, _ := db.Prepare("UPDATE personagens set nome = ? where id = ?")
	stmt.Exec(nome, id)
	stmt.Close()
}
