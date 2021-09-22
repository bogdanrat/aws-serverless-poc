package publisher

type Publisher interface {
	Publish(string) error
}
