package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/scenery/mediax/handlers"
	"github.com/scenery/mediax/helpers"
	"github.com/scenery/mediax/models"
)

func redirectToHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
}

func handleHomePage(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	weekdays := [...]string{"星期日", "星期一", "星期二", "星期三", "星期四", "星期五", "星期六"}
	today := fmt.Sprintf("%d月%d日 %s", now.Month(), now.Day(), weekdays[now.Weekday()])

	recentSubjects, err := handlers.GetRecentSubjects(5)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get recent subjects: %v", err), http.StatusInternalServerError)
		return
	}

	var recentGroups []models.HomeViewType
	subjectTypes := []string{"book", "movie", "tv", "anime", "game"}
	for _, subjectType := range subjectTypes {
		summary, _ := handlers.GetHomeSummary(subjectType)
		recentGroups = append(recentGroups, models.HomeViewType{
			SubjectType: subjectType,
			TypeZH:      helpers.GetTypeZH(subjectType),
			ActionZH:    helpers.GetActionZH(subjectType),
			UnitZH:      helpers.GetUnitZH(subjectType),
			Items:       processHomeHTML(recentSubjects[subjectType]),
			Summary:     summary,
		})
	}

	data := models.HomeView{
		Today:        today,
		PageTitle:    "主页",
		RecentGroups: recentGroups,
	}

	renderPage(w, "index.html", data)
}
