package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rikuya98/go-poke-data-api/models"
	"github.com/rikuya98/go-poke-data-api/services"
	"net/http"
	"net/url"
	"strconv"
	"sync"
)

func GetQueryParams(que url.Values, keys []string) (models.PokeParams, error) {
	var pokeQueParams models.PokeParams
	var err error

	for _, key := range keys {
		val := que.Get(key)
		if val == "" {
			return models.PokeParams{}, ErrNoQuery
		}
		switch key {
		case "lv":
			pokeQueParams.Level, err = strconv.Atoi(val)
		case "ef":
			pokeQueParams.EffortVal, err = strconv.Atoi(val)
		case "in":
			pokeQueParams.IndVal, err = strconv.Atoi(val)
		default:
			return models.PokeParams{}, ErrInvalidKey
		}
		if err != nil {
			return models.PokeParams{}, err
		}
	}
	return pokeQueParams, nil
}

// GetPokeDataHandler ポケモンのデータを取得する。
func GetPokeDataHandler(w http.ResponseWriter, req *http.Request) {
	var err error
	var getImgErr error
	var getDataErr error
	var getNameErr error

	var pokeParams models.PokeParams
	var pokeData models.PokeData
	var wg sync.WaitGroup
	//クエリパラメータからレベル、努力値、個体値を取得
	query := req.URL.Query()
	keys := []string{"lv", "ef", "in"}
	pokeParams, err = GetQueryParams(query, keys)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//パスパラメータを元にポケモンのデータを取得
	vars := mux.Vars(req)
	id := vars["id"]

	wg.Add(3)

	//ポケモンの画像を取得
	var pokeEncData string
	go func() {
		defer wg.Done()
		pokeEncData, getImgErr = services.GetPokeImageService(id)
	}()

	var pokeName string
	go func() {
		defer wg.Done()
		pokeName, getNameErr = services.GetPokeNameService(id)
	}()

	go func() {
		defer wg.Done()
		pokeData, getDataErr = services.GetPokeDataService(id)
	}()

	wg.Wait()

	if getNameErr != nil {
		http.Error(w, getNameErr.Error(), http.StatusInternalServerError)
		return
	}

	if getImgErr != nil {
		http.Error(w, getImgErr.Error(), http.StatusInternalServerError)
		return
	}

	if getDataErr != nil {
		http.Error(w, getDataErr.Error(), http.StatusInternalServerError)
		return
	}

	pokeData.EncImg = pokeEncData
	pokeData.Name = pokeName
	// レベル、努力値、個体値を元にステータスを計算
	services.CalPokeStat(&pokeData, pokeParams)

	//レスポンスをjson形式で返す
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(pokeData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
