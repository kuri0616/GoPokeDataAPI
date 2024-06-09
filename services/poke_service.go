package services

import (
	"encoding/base64"
	"encoding/json"
	"github.com/rikuya98/go-poke-data-api/models"
	"io"
	"net/http"
)

// CalHP ポケモンのHPを計算する
func CalHP(baseStat int, individualVal int, effortVal int, level int) int {
	calStatus := (baseStat*2+individualVal+effortVal/4)*level/100 + level + 10
	return calStatus
}

// CalOtherStat ポケモンのHP以外のステータスを計算する
func CalOtherStat(baseStat int, individualVal int, effortVal int, level int) int {
	calStatus := (baseStat*2+individualVal+effortVal/4)*level/100 + 5
	return calStatus
}

// GetPokeImageService ポケモンの画像を取得する
func GetPokeImageService(id string) (encImg string, err error) {
	//ポケモンの画像を取得
	var res *http.Response

	res, err = http.Get("https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/" + id + ".png")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	pokeImg, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(pokeImg), nil
}

// GetPokeNameService ポケモンの名前を取得する
func GetPokeNameService(id string) (name string, err error) {
	var findName models.FindName
	var jaName string

	resName, err := http.Get("https://pokeapi.co/api/v2/pokemon-species/" + id)
	if err != nil {
		return "", err
	}
	defer resName.Body.Close()

	if err := json.NewDecoder(resName.Body).Decode(&findName); err != nil {
		return "", err
	}

	//日本語名を取得
	for _, name := range findName.Names {
		if name.Language.Name == "ja" {
			jaName = name.Name
			break
		}
	}
	return jaName, nil
}

// GetPokeDataService ポケモンのデータを取得する
func GetPokeDataService(id string) (pokeData models.PokeData, err error) {
	var PokeData models.PokeData
	res, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + id)

	if err != nil {
		return PokeData, err
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&PokeData); err != nil {
		return PokeData, err
	}
	return PokeData, nil
}

// CalPokeStat ポケモンのステータスを計算する
func CalPokeStat(pokeData *models.PokeData, params models.PokeParams) {
	// レベル、努力値、個体値を元にステータスを計算
	for i, stat := range pokeData.Stats {
		switch stat.Stat.Name {
		case "hp":
			pokeData.Stats[i].CalStat = CalHP(stat.BaseStat, params.IndVal, params.EffortVal, params.Level)
		default:
			pokeData.Stats[i].CalStat = CalOtherStat(stat.BaseStat, params.IndVal, params.EffortVal, params.Level)
		}
	}
}
