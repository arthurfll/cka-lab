package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "modernc.org/sqlite"
)

type Registro struct {
	Valor int `json:"valor"`
}

func main() {
	// Conecta ao banco SQLite ou cria o arquivo app.db.
	db, err := sql.Open("sqlite", "/app/data/app.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	// Cria a tabela caso ela não exista.
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS registros (
			id INTEGER PRIMARY KEY CHECK (id = 1),
			valor INTEGER NOT NULL
		);
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Cria o registro inicial com valor 0, somente se não existir.
	_, err = db.Exec(`
		INSERT OR IGNORE INTO registros (id, valor)
		VALUES (1, 0);
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Endpoint único.
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		// Inicia uma transação.
		tx, err := db.Begin()
		if err != nil {
			http.Error(w, "Erro ao iniciar transação", http.StatusInternalServerError)
			return
		}
		defer tx.Rollback()

		// Incrementa o registro existente.
		_, err = tx.Exec(`
			UPDATE registros
			SET valor = valor + 1
			WHERE id = 1;
		`)
		if err != nil {
			http.Error(w, "Erro ao incrementar registro", http.StatusInternalServerError)
			return
		}

		// Consulta o valor atualizado.
		var registro Registro

		err = tx.QueryRow(`
			SELECT valor
			FROM registros
			WHERE id = 1;
		`).Scan(&registro.Valor)

		if err != nil {
			http.Error(w, "Erro ao consultar registro", http.StatusInternalServerError)
			return
		}

		// Confirma a transação.
		if err := tx.Commit(); err != nil {
			http.Error(w, "Erro ao salvar registro", http.StatusInternalServerError)
			return
		}

		// Retorna erro 500 se o valor ultrapassar 5.
		if registro.Valor > 5 {
			http.Error(w, "Limite excedido", http.StatusInternalServerError)
			return
		}

		// Exibe o log no terminal.
		log.Printf("Registro iterado: valor = %d", registro.Valor)

		// Retorna o valor atualizado em JSON.
		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(registro); err != nil {
			log.Printf("Erro ao enviar resposta: %v", err)
		}
	})

	log.Println("API disponível em http://localhost:7000")

	log.Fatal(http.ListenAndServe(":7000", nil))
}