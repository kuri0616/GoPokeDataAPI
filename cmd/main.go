package main

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strconv"
)

func main() {
	//ポケモンのデータを格納する構造体
	type PokeData struct {
		EncImg string `json:"img"`
		Stats  []struct {
			BaseStat int `json:"base_stat"`
			CalStat  int `json:"cal_stat"`
			Stat     struct {
				Name string `json:"name"`
			} `json:"stat"`
		} `json:"stats"`
	}

	//最大HP
	//(種族値×2+個体値+努力値÷4)×レベル÷100+レベル+10
	//こうげき・ぼうぎょ・とくこう・とくぼう・すばやさ
	//{(種族値×2+個体値+努力値÷4)×レベル÷100+5}×せいかく補正
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
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			effortVal, err = strconv.Atoi(pathEffortVal)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			indvidualVal, err = strconv.Atoi(pathIndividualVal)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		vars := mux.Vars(req)
		//パスパラメータを元にポケモンのデータを取得
		res, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + vars["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer res.Body.Close()

		//ポケモンの画像を取得
		resImg, err := http.Get("https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/" + vars["id"] + ".png")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resImg.Body.Close()

		pokeImg, err := io.ReadAll(resImg.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pokeEncData := base64.StdEncoding.EncodeToString(pokeImg)

		var PokeData PokeData
		//レスポンスを構造体に変換
		if err := json.NewDecoder(res.Body).Decode(&PokeData); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		PokeData.EncImg = pokeEncData

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
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
