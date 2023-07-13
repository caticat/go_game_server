package phelp

type PSortedMap[K comparable, V any] struct {
	M_mapData map[K]V `json:"data"`
	M_sliKey  []K     `json:"-"`
}

func NewPSortedMap[K comparable, V any]() *PSortedMap[K, V] {
	return &PSortedMap[K, V]{
		M_mapData: make(map[K]V),
		M_sliKey:  make([]K, 0),
	}
}

func (t *PSortedMap[K, V]) InitByMap(mapData map[K]V) {
	t.Clear()

	for k, v := range mapData {
		t.Set(k, v)
	}
}

func (t *PSortedMap[K, V]) GetMap() map[K]V {
	return t.M_mapData
}

func (t *PSortedMap[K, V]) Length() int {
	return len(t.M_sliKey)
}

func (t *PSortedMap[K, V]) GetByIndex(index int) (key K, value V, err error) {
	l := t.Length()
	if index >= l {
		err = ErrorIndexOutofRange
		return
	}

	key = t.M_sliKey[index]
	value = t.M_mapData[key]
	return
}

func (t *PSortedMap[K, V]) Get(key K) (V, bool) {
	value, ok := t.M_mapData[key]
	return value, ok
}

func (t *PSortedMap[K, V]) Set(key K, value V) {
	if _, ok := t.M_mapData[key]; !ok {
		t.M_sliKey = append(t.M_sliKey, key)
	}

	t.M_mapData[key] = value
}

func (t *PSortedMap[K, V]) Del(key K) {
	if _, ok := t.M_mapData[key]; !ok {
		return
	}

	for i, k := range t.M_sliKey {
		if k != key {
			continue
		}
		t.M_sliKey = append(t.M_sliKey[:i], t.M_sliKey[i+1:]...)
		break
	}

	delete(t.M_mapData, key)
}

func (t *PSortedMap[K, V]) Clear() {
	t.M_sliKey = make([]K, 0)

	for k := range t.M_mapData {
		delete(t.M_mapData, k)
	}
}

func (t *PSortedMap[K, V]) Append(mapData map[K]V) {
	for k, v := range mapData {
		t.Set(k, v)
	}
}
