package u2semi

type RequestRepository interface {
	Save(req *Request) error

	Migrate() error
}
