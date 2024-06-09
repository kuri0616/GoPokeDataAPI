package handlers

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/rikuya98/go-poke-data-api/models"
	"github.com/rikuya98/go-poke-data-api/services"
	"net/http"
	"net/url"
	"strconv"
)

var ErrNoQuery = errors.New("no query parameter")

// ステータスを計算する関数

func CalHP(baseStat int, individualVal int, effortVal int, level int) int {
	calStatus := (baseStat*2+individualVal+effortVal/4)*level/100 + level + 10
	return calStatus
}
func CalOtherStat(baseStat int, individualVal int, effortVal int, level int) int {
	calStatus := (baseStat*2+individualVal+effortVal/4)*level/100 + 5
	return calStatus
}

func GetQueryIntParams(que url.Values, key string) (int, error) {
	val := que.Get(key)
	if val == "" {
		return 0, ErrNoQuery
	}
	return strconv.Atoi(val)
}

//ポケモンのデータを取得するhandler

func GetPokeDataHandler(w http.ResponseWriter, req *http.Request) {
	var err error
	var level, effortVal, indvidualVal int

	//クエリパラメータからレベル、努力値、個体値を取得
	query := req.URL.Query()

	level, err = GetQueryIntParams(query, "lv")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	effortVal, err = GetQueryIntParams(query, "ef")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	indvidualVal, err = GetQueryIntParams(query, "in")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//パスパラメータを元にポケモンのデータを取得
	vars := mux.Vars(req)
	id := vars["id"]
	res, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	var FindName models.FindName
	//ポケモンの名前を取得
	resName, err := http.Get("https://pokeapi.co/api/v2/pokemon-species/" + id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resName.Body.Close()

	if err := json.NewDecoder(resName.Body).Decode(&FindName); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//ポケモンの画像を取得

	pokeEncData, err := services.GetPokeImageService(id)

	//レスポンスを構造体に変換
	var PokeData models.PokeData
	if err := json.NewDecoder(res.Body).Decode(&PokeData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//日本語の名前を取得
	for _, name := range FindName.Names {
		if name.Language.Name == "ja" {
			PokeData.Name = name.Name
			break
		}
	}

	PokeData.EncImg = pokeEncData

	// レベル、努力値、個体値を元にステータスを計算
	for i, stat := range PokeData.Stats {
		switch stat.Stat.Name {
		case "hp":
			PokeData.Stats[i].CalStat = CalHP(stat.BaseStat, indvidualVal, effortVal, level)
		default:
			PokeData.Stats[i].CalStat = CalOtherStat(stat.BaseStat, indvidualVal, effortVal, level)
		}
	}

	//レスポンスをjson形式で返す
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(PokeData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
