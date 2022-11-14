package rest

type Rest interface {
	GetGin() *Gin
	Run()
}
