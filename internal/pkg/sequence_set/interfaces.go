package set

type Set interface {
	Add(key string, sequenceNumber int) error
	Has(key string) (bool, error)
	Remove(key string, sequenceNumber int) error
}
