package config

import "time"

const (
	// Image
	ImageDir = "images"

	// Cache
	MaxCacheSubjects = 1000
)

// API Config
var (
	CORS_HOST  = "*"
	QueryLimit = 50
)

var DefaultConfig = AppConfig{
	Server: ServerConfig{
		Address: "0.0.0.0",
		Port:    8080,
	},
	SessionTimeout: 7 * 24 * time.Hour,
	Pagination: PaginationConfig{
		PageSize: 10,
	},
	Categories: defaultCategories,
}

var defaultCategories = []string{"book", "movie", "tv", "anime", "game"}

type CategoryInfo struct {
	Name        string
	Unit        string
	ActionFull  string
	ActionShort string
}

var CategoryInfoMap = map[string]CategoryInfo{
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
