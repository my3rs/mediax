package routes

import (
	"fmt"
	"html/template"
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
	hour := time.Now().Hour()
	var greeting string
	switch {
	case hour >= 5 && hour < 9:
		greeting = "早上好"
	case hour >= 9 && hour < 12:
		greeting = "上午好"
	case hour >= 12 && hour < 14:
		greeting = "中午好"
	case hour >= 14 && hour < 18:
		greeting = "下午好"
	case hour >= 18 && hour < 24:
		greeting = "晚上好"
	default:
		greeting = "深夜好"
	}

	recentSubjects, err := handlers.GetRecentSubjects(5)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get recent subjects: %v", err), http.StatusInternalServerError)
		return
	}

	var recentGroups []models.HomeViewType
	subjectTypes := helpers.GetCategories()
	for _, subjectType := range subjectTypes {
		subjectActionFullName, subjectActionShortName := helpers.GetSubjectActionName(subjectType)
		summary, _ := handlers.GetHomeSummary(subjectType)
		recentGroups = append(recentGroups, models.HomeViewType{
			SubjectType:            subjectType,
			SubjectTypeName:        helpers.GetSubjectTypeName(subjectType),
			SubjectActionFullName:  subjectActionFullName,
			SubjectActionShortName: subjectActionShortName,
			SubjectUnitName:        helpers.GetSubjectUnitName(subjectType),
			CategoryIcon:           template.HTML(helpers.GetCategoryIcon(subjectType, "30", fmt.Sprintf("var(--%s-color)", subjectType))),
			Items:                  processHomeHTML(recentSubjects[subjectType]),
			Summary:                summary,
		})
	}

	data := models.HomeView{
		Header:       helpers.GetHeader("home"),
		TimePeriod:   greeting,
		PageTitle:    "主页",
		RecentGroups: recentGroups,
	}

	renderPage(w, "index.html", data)
}
