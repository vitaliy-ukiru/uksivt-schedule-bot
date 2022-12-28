package group

type Usecase interface {
	Years() []int
	Specs(year int) []string
	Numbers(year int, spec string) []int
}
