package web

type Web interface {
	Client() *Gin
	Run() error
}
