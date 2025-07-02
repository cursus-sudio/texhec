package audio

type Api interface {
}

type api struct{}

func newApi() Api {
	return api{}
}
