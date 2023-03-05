package model

type Student struct {
	StudentID string `json:"student_id"`
	Name      string `json:"name"`
	College   string `json:"college"`
	Major     string `json:"major"`
	Class     string `json:"class"`
}

type Score struct {
	Term      string `json:"term"`
	ClassName string `json:"class_name"`
	Score     string `json:"score"`
	GPA       string `json:"gpa"`
	Pattern   string `json:"pattern"`
	Credits   string `json:"credits"`
}

type ScoreResult struct {
	ScoreList          []Score `json:"score_list"`
	CourseCount        string  `json:"course_count"`
	TotalCredit        string  `json:"total_credit"`
	AverageCreditPoint string  `json:"average_credit_point"`
	AverageGrade       string  `json:"average_grade"`
	Rank               string  `json:"rank"`
}

type Course struct {
	CourseName string `json:"course_name"`
	WeekList   []int  `json:"week_list"`
	Week       int    `json:"week"`
	Address    string `json:"address"`
	Teacher    string `json:"teacher"`
	Sections   int    `json:"sections"`
}

type CourseDetail struct {
	Type        string `json:"type"`
	CourseName  string `json:"course_name"`
	Status      string `json:"status"`
	Attribute   string `json:"attribute"`
	Credit      string `json:"credit"`
	CreditHours string `json:"credit_hours"`
	Term        string `json:"term"`
}
