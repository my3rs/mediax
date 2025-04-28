package config

import "github.com/scenery/mediax/models"

const (
	// HTTP
	HTTP_PORT = 8080

	// Image
	ImageDir = "images"

	// Page
	PageSize = 10

	// Cache
	MaxCacheSubjects = 1000
)

// API Config
var (
	CORS_HOST    = "*"
	RequestLimit = 50
)

// Categories
var (
	Categories      = []string{"book", "movie", "tv", "anime", "game"}
	CategoryInfoMap = map[string]models.CategoryInfo{
		"book": {
			Name:        "图书",
			Unit:        "本",
			ActionFull:  "阅读",
			ActionShort: "读",
		},
		"movie": {
			Name:        "电影",
			Unit:        "部",
			ActionFull:  "观看",
			ActionShort: "看",
		},
		"tv": {
			Name:        "剧集",
			Unit:        "部",
			ActionFull:  "观看",
			ActionShort: "看",
		},
		"anime": {
			Name:        "番剧",
			Unit:        "部",
			ActionFull:  "观看",
			ActionShort: "看",
		},
		"game": {
			Name:        "游戏",
			Unit:        "款",
			ActionFull:  "游玩",
			ActionShort: "玩",
		},
		"home": {
			Name: "mediaX",
		},
		"search": {
			Name: "搜索",
		},
	}
)
