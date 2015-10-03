package center

import ()

type mapHandler interface {
	OutputMapByJson(json []byte) (map[string]string, error)
}
type mapDataHandler interface {
	StoreMapData(data map[string]string) error
}
type Server interface {
	Start() error
	GetDataChannel() chan [config.MaxNumberOfData][]byte
}
