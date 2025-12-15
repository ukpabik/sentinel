package main

type Limiter interface {
	Allow() bool
}
