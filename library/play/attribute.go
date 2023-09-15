package play

import "sync"

type AttributesKey string

func (ak AttributesKey) String() string {
	return string(ak)
}

type Attributes struct {
	_map sync.Map
}

func (att *Attributes) Put(key AttributesKey, value interface{}) {
	//att.Lock()
	//att.Map[key] = value
	att._map.Store(key, value)
	//defer att.Unlock()
}

func (att *Attributes) GetOrPut(key AttributesKey, value interface{}) (actual interface{}, loaded bool) {
	//att.Lock()
	//att.Map[key] = value
	return att._map.LoadOrStore(key, value)
	//defer att.Unlock()
}
func (att *Attributes) GetMap() map[string]interface{} {
	data := make(map[string]interface{})

	att._map.Range(func(key, value interface{}) bool {

		data[key.(AttributesKey).String()] = value

		return true
	})
	return data
}
func (att *Attributes) Get(key AttributesKey) interface{} {
	//att.RLock()
	//defer att.RUnlock()
	v, _ := att._map.Load(key)
	return v
}
func (att *Attributes) Delete(key AttributesKey) {
	//att.RLock()
	//defer att.RUnlock()
	//delete(att.Map, key)
	att._map.Delete(key)
}
