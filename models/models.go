package models

// PokeData ポケモンのデータを格納する構造体
type PokeData struct {
	Name   string `json:"name"`
	EncImg string `json:"img"`
	Stats  []struct {
		BaseStat int `json:"base_stat"`
		CalStat  int `json:"cal_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
}

// FindName ポケモンの日本語名を取得するための構造体
type FindName struct {
	Names []struct {
		Language struct {
			Name string `json:"name"`
		} `json:"language"`
		Name string `json:"name"`
	}
}

// PokeParams ポケモンのステータスを計算するためのパラメータ
type PokeParams struct {
	IndVal    int
	EffortVal int
	Level     int
}
