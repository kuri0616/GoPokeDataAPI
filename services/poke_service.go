package services

import (
	"encoding/base64"
	"io"
	"net/http"
)

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
