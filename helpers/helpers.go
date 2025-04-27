package helpers

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func GenerateUUID() string {
	return uuid.New().String()
}

func GetTimestamp() int64 {
	return time.Now().Unix()
}

func StringToInt(value string) (int, error) {
	result, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func GetTypeZH(subjectType string) string {
	switch subjectType {
	case "book":
		return "图书"
	case "movie":
		return "电影"
	case "tv":
		return "剧集"
	case "anime":
		return "番剧"
	case "game":
		return "游戏"
	default:
		return "未知"
	}
}

func GetActionZH(subjectType string) string {
	switch subjectType {
	case "book":
		return "阅读"
	case "movie", "tv", "anime":
		return "观看"
	case "game":
		return "游玩"
	default:
		return "操作"
	}
}

func GetUnitZH(subjectType string) string {
	switch subjectType {
	case "book":
		return "本"
	case "movie", "tv", "anime":
		return "部"
	case "game":
		return "款"
	default:
		return ""
	}
}

func MD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return fmt.Sprintf("%x", hash)
}
