package sentinel

type Limiter interface {
	Allow() bool
}
