package helpers

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/scenery/mediax/config"
	"github.com/scenery/mediax/models"
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

func GetHeader(currentCategory string) models.Header {
	var options []models.HeaderOption
	for _, cat := range GetCategories() {
		options = append(options, models.HeaderOption{
			Category:     cat,
			CategoryName: GetSubjectTypeName(cat),
		})
	}

	return models.Header{
		Options:     options,
		Current:     currentCategory,
		CurrentName: GetSubjectTypeName(currentCategory),
	}
}

func GetCategories() []string {
	return append([]string{}, config.Categories...)
}

func GetSubjectTypeName(subjectType string) string {
	if info, ok := config.CategoryInfoMap[subjectType]; ok {
		return info.Name
	}
	return "未知"
}

func GetSubjectActionName(subjectType string) (string, string) {
	if info, ok := config.CategoryInfoMap[subjectType]; ok {
		return info.ActionFull, info.ActionShort
	}
	return "", ""
}

func GetSubjectUnitName(subjectType string) string {
	if info, ok := config.CategoryInfoMap[subjectType]; ok {
		return info.Unit
	}
	return ""
}

func MD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return fmt.Sprintf("%x", hash)
}
