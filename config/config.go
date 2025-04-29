package config

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/scenery/mediax/models"
)

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

var DefaultConfig = models.AppConfig{
	Server: models.ServerConfig{
		Address: "0.0.0.0",
		Port:    8080,
	},
	Pagination: models.PaginationConfig{
		PageSize: 10,
	},
	Categories: defaultCategories,
}

var defaultCategories = []string{"book", "movie", "tv", "anime", "game"}

var CategoryInfoMap = map[string]models.CategoryInfo{
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

var App models.AppConfig

func LoadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg models.AppConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("invalid config format: %w", err)
	}

	if cfg.Server.Address == "" {
		cfg.Server.Address = DefaultConfig.Server.Address
	} else if net.ParseIP(cfg.Server.Address) == nil {
		return fmt.Errorf("invalid server address: %s", cfg.Server.Address)
	}

	if cfg.Server.Port == 0 {
		cfg.Server.Port = DefaultConfig.Server.Port
	} else if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
		return fmt.Errorf("invalid port: %d (must be between 1 and 65535)", cfg.Server.Port)
	}

	if cfg.Pagination.PageSize == 0 {
		cfg.Pagination.PageSize = DefaultConfig.Pagination.PageSize
	} else if cfg.Pagination.PageSize < 1 || cfg.Pagination.PageSize > 50 {
		return fmt.Errorf("invalid pagination page_size: %d (must be between 1 and 50)", cfg.Pagination.PageSize)
	}

	if len(cfg.Categories) == 0 {
		cfg.Categories = DefaultConfig.Categories
	} else {
		for _, cat := range cfg.Categories {
			if !isValidCategory(cat) {
				return fmt.Errorf("invalid category: %s (allowed: %v)", cat, defaultCategories)
			}
		}
	}

	App = cfg
	return nil
}

func isValidCategory(cat string) bool {
	for _, valid := range defaultCategories {
		if cat == valid {
			return true
		}
	}
	return false
}
