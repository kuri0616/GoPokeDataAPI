package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func main() {
	//ポケモンのデータを格納する構造体
	type PokeData struct {
		Stats []struct {
			BaseStat int `json:"base_stat"`
			CalStat  int `json:"cal_stat"`
			Stat     struct {
				Name string `json:"name"`
			} `json:"stat"`
		} `json:"stats"`
	}

	//ステータスを計算する関数
	CalculateHP := func(baseStat int, individualVal int, effortVal int, level int) int {
		calStatus := (baseStat*2+individualVal+effortVal/4)*level/100 + level + 10
		return calStatus
	}
	CalculateOtherStat := func(baseStat int, individualVal int, effortVal int, level int) int {
		calStatus := (baseStat*2+individualVal+effortVal/4)*level/100 + 5
		return calStatus
	}

	//ポケモンのデータを取得するhandler
	GetPokeDataHandler := func(w http.ResponseWriter, req *http.Request) {
		var err error
		var level, effortVal, indvidualVal int

		//クエリパラメータからレベル、努力値、個体値を取得
		query := req.URL.Query()
		pathLevel := query.Get("lv")
		pathEffortVal := query.Get("ef")
		pathIndividualVal := query.Get("in")

		if pathLevel != "" && pathEffortVal != "" && pathIndividualVal != "" {
			level, err = strconv.Atoi(pathLevel)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			effortVal, err = strconv.Atoi(pathEffortVal)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			indvidualVal, err = strconv.Atoi(pathIndividualVal)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		//パスパラメータを元にポケモンのデータを取得
		vars := mux.Vars(req)
		res, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + vars["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer res.Body.Close()

		//レスポンスを構造体に変換
		var PokeData PokeData
		if err := json.NewDecoder(res.Body).Decode(&PokeData); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// レベル、努力値、個体値を元にステータスを計算
		for i, stat := range PokeData.Stats {
			switch stat.Stat.Name {
			case "hp":
				PokeData.Stats[i].CalStat = CalculateHP(stat.BaseStat, indvidualVal, effortVal, level)
			default:
				PokeData.Stats[i].CalStat = CalculateOtherStat(stat.BaseStat, indvidualVal, effortVal, level)
			}
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
	r.HandleFunc("/pokemon/{id:[0-9]+}", GetPokeDataHandler).Methods(http.MethodGet)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
