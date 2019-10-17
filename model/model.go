package model

type (
	// Song 单曲
	Song struct {
		ID         int      `json:"id"`
		Mid        string   `json:"mid"`
		Title      string   `json:"title"`
		Singer     []Singer `json:"singer"`
		Album      Album    `json:"album"`
		IndexAlbum int      `json:"index_album"`
		TimePublic string   `json:"time_public"`
	}

	// Singer 歌手
	Singer struct {
		ID   int    `json:"id"`
		Mid  string `json:"mid"`
		Name string `json:"name"`
	}

	// Album 专辑
	Album struct {
		ID   int    `json:"id"`
		Mid  string `json:"mid"`
		Name string `json:"name"`
	}
)
