package bot

type Bot interface {
	Run()
	Init() error
}
