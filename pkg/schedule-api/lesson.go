package scheduleapi

type Lesson struct {
	Group        string `json:"college_group"`
	DayOfWeek    int    `json:"day_of_week"`
	LessonNumber int    `json:"lesson_number"`
	Name         string `json:"lesson"`
	Teacher      string `json:"teacher"`
	LessonHall   string `json:"lesson_hall"`
	Time         string `json:"time"`
	Replacement  bool   `json:"replacement"`
}
