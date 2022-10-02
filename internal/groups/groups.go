package groups

type Service interface {
	Years() []int
	Specs(year int) []string
	Numbers(year int, spec string) []int
}
