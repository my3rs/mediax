package models

import "html/template"

type SubjectView struct {
	Category        string
	PageTitle       string
	ManageType      int
	CreatorLabel    string
	PressLabel      string
	PubDateLabel    string
	SummaryLabel    string
	StatusText      string
	RatingStar      int
	ImageURL        string
	ExternalURLIcon template.HTML
	Subject         Subject
}

type CategoryViewItem struct {
	SubjectType  string
	SubjectURL   string
	Title        string
	AltTitle     string
	Creator      string
	Press        string
	PubDate      string
	MarkDate     string
	Rating       int
	StatusText   string
	CreatorLabel string
	PressLabel   string
	PubDateLabel string
	ImageURL     string
}

type StatusCounts struct {
	All     int64
	Todo    int64
	Doing   int64
	Done    int64
	OnHold  int64
	Dropped int64
}

type StatusLabel struct {
	Value int
	Label string
	Count int64
}

type CategoryView struct {
	Category    string
	PageTitle   string
	Status      int
	TotalCounts int64
	StatusList  []StatusLabel
	SortBy      int
	CurrentPage int
	TotalPages  int
	PageNumbers []int
	PageParams  template.URL
	Subjects    []CategoryViewItem
}

type SearchView struct {
	Category    string
	PageTitle   string
	Query       string
	QueryType   string
	TotalCount  int
	CurrentPage int
	TotalPages  int
	PageNumbers []int
	PageParams  template.URL
	Subjects    []CategoryViewItem
}

type AddOptions struct {
	SubjectType string
	TypeZH      string
}

type AddView struct {
	PageTitle  string
	Category   string
	AddType    string
	AddOptions []AddOptions
}

type HomeLastItem struct {
	Title      string
	SubjectURL string
	Status     int
	Date       string
}

type HomeSummary struct {
	MonthCount    int
	HalfYearCount int
	YearCount     int
	LastItem      HomeLastItem
}

type HomeViewItem struct {
	SubjectURL string
	ImageURL   string
}

type HomeViewType struct {
	SubjectType string
	TypeZH      string
	ActionZH    string
	UnitZH      string
	Items       []HomeViewItem
	Summary     HomeSummary
}

type HomeView struct {
	Category     string
	Today        string
	PageTitle    string
	RecentGroups []HomeViewType
}
