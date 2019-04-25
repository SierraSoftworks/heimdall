package file

type FileFormat struct {
	Channel string `json:"channel"`
	Data    []byte `json:"data"`
}
