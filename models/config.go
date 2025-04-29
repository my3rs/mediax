package models

type ServerConfig struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
}

type PaginationConfig struct {
	PageSize int `json:"page_size"`
}

type AppConfig struct {
	Server     ServerConfig     `json:"server"`
	Pagination PaginationConfig `json:"pagination"`
	Categories []string         `json:"categories"`
}

type CategoryInfo struct {
	Name        string
	Unit        string
	ActionFull  string
	ActionShort string
}
