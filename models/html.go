package models

import "html/template"

type HeaderOption struct {
	Category     string
	CategoryName string
}

type Header struct {
	Options     []HeaderOption
	Current     string
	CurrentName string
}

type CategoryInfo struct {
	Name        string
	Unit        string
	ActionFull  string
	ActionShort string
}

type SubjectView struct {
	Header          Header
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
	Header      Header
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
	Header      Header
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

type AddView struct {
	Header    Header
	PageTitle string
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
	SubjectType            string
	SubjectTypeName        string
	SubjectActionFullName  string
	SubjectActionShortName string
	SubjectUnitName        string
	Items                  []HomeViewItem
	Summary                HomeSummary
}

type HomeView struct {
	Header       Header
	Today        string
	PageTitle    string
	RecentGroups []HomeViewType
}
