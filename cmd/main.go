package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	type Stat struct {
		Name string `json:"name"`
	}
	type Stats struct {
		BaseStat int  `json:"base_stat"`
		Stat     Stat `json:"stat"`
	}
	type PokeData struct {
		Stats []Stats `json:"stats"`
	}

	//ポケモンのデータを取得するhandler
	GetPokeDataHandler := func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		//パスパラメータを数値に変換
		var PokeData PokeData
		res, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + vars["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer res.Body.Close()

		//レスポンスを構造体に変換
		if err := json.NewDecoder(res.Body).Decode(&PokeData); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//レスポンスをjson形式で返す
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(PokeData); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	r := mux.NewRouter()
	r.HandleFunc("/pokemon/{id:[0-9]+}", GetPokeDataHandler).Methods(http.MethodGet)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
