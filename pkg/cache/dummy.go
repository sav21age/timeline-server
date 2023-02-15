package cache

type Dummy struct {}

func NewDummy() *Dummy {
	return &Dummy{}
}

func (c *Dummy) Set(key, value interface{}, ttl int64) error {
	return nil
}

func (c *Dummy) Get(key interface{}) (interface{}, error){
	return key, ErrItemNotFound
}
