package redis

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/tool/object"
	"reflect"
	"sync"
	"time"
)

type local struct {
	data sync.Map
}

type MapItem struct {
	V        any
	TTL      time.Duration
	CreateAt time.Time
}

func (l *local) Del(ctx context.Context, keys ...string) (int64, error) {
	for k := range keys {
		l.data.Delete(k)
	}
	return int64(len(keys)), nil
}

func (l *local) Get(ctx context.Context, key string) (string, error) {
	value, ok := l.data.Load(key)
	if !ok {
		return "", errors.New("not exist")
	}
	item := value.(*MapItem)
	if time.Now().Sub(item.CreateAt) >= item.TTL {
		_, err := l.Del(ctx, key)
		if err != nil {
			return "", err
		}
		return "", errors.New("not exist")
	}
	return object.ParseString(item.V), nil
}

func (l *local) GetEx(ctx context.Context, key string, expiration time.Duration) (string, error) {
	value, ok := l.data.Load(key)
	if !ok {
		return "", errors.New("not exist")
	}
	item := value.(*MapItem)
	if time.Now().Sub(item.CreateAt) >= item.TTL {
		_, err := l.Del(ctx, key)
		if err != nil {
			return "", err
		}
		return "", errors.New("not exist")
	}
	item.TTL = expiration
	item.CreateAt = time.Now()
	l.data.Store(key, item)
	return object.ParseString(item.V), nil
}

func (l *local) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	itemValue := ""
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		if v.Kind() == reflect.Struct || v.Kind() == reflect.Map || v.Kind() == reflect.Slice {
			marshal, err := json.Marshal(value)
			if err != nil {
				return err
			}
			itemValue = string(marshal)
		}
	} else {
		itemValue = object.ParseString(value)
	}
	item := &MapItem{}
	item.V = itemValue
	item.TTL = expiration
	item.CreateAt = time.Now()
	l.data.Store(key, item)
	return nil
}

func (l *local) TryLock(ctx context.Context, key string, wait ...time.Duration) (bool, func()) {
	return true, func() {

	}
}

func (l *local) GenerateUID(ctx context.Context, maxID int64) (uint64, error) {
	key := NewUIDKey()
	mUID, _ := l.Get(ctx, key)
	v := object.ParseUint(mUID)
	if v == 0 {
		if maxID < 1000000 {
			maxID = 1000001
		}
		var err error
		_, err = l.Incr(ctx, key)
		if err != nil {
			return 0, err
		}
	}
	vv, err := l.Incr(ctx, key)
	if err != nil {
		return 0, err
	}
	return uint64(vv), nil
}

func (l *local) Expire(ctx context.Context, key string, expiration time.Duration) error {
	value, ok := l.data.Load(key)
	if !ok {
		return errors.New("not exist")
	}
	item := value.(*MapItem)
	item.TTL = expiration
	item.CreateAt = time.Now()
	l.data.Store(key, item)
	return nil
}

func (l *local) HSet(ctx context.Context, key string, value map[string]any) error {
	item := &MapItem{}
	item.V = value
	item.TTL = 0
	item.CreateAt = time.Now()
	l.data.Store(key, item)
	return nil
}

func (l *local) HMGet(ctx context.Context, key string, fields ...string) ([]any, error) {
	value, ok := l.data.Load(key)
	if !ok {
		return nil, errors.New("not exist")
	}
	var arr []any
	vm := value.(map[string]any)
	for i := range fields {
		v, ok := vm[fields[i]]
		if ok {
			arr = append(arr, v)
		}
	}
	return arr, nil
}

func (l *local) HGet(ctx context.Context, key, field string) (string, error) {
	value, ok := l.data.Load(key)
	if !ok {
		return "", errors.New("not exist")
	}
	v, ok := value.(map[string]any)[field]
	if !ok {
		return "", errors.New("not exist key")
	}

	return object.ParseString(v), nil
}

func (l *local) Exists(ctx context.Context, keys ...string) (int64, error) {
	for _, key := range keys {
		_, ok := l.data.Load(key)
		if !ok {
			return 0, errors.New("not exist")
		}
	}
	return int64(len(keys)), nil
}

func (l *local) Incr(ctx context.Context, key string) (int64, error) {
	value, ok := l.data.Load(key)
	if !ok {
		value = MapItem{V: "0", CreateAt: time.Now(), TTL: 0}
	}
	item := value.(*MapItem)
	v := object.ParseInt(value) + 1
	item.V = v
	l.data.Store(key, item)
	return int64(v), nil
}

func (l *local) SetAdd(ctx context.Context, key string, members ...any) (int64, error) {
	value, ok := l.data.Load(key)
	if !ok {
		value = MapItem{V: []any{}, CreateAt: time.Now(), TTL: 0}
	}
	item := value.(*MapItem)
	vs := item.V.([]any)
	vs = append(vs, members...)
	item.V = vs
	l.data.Store(key, item)
	return int64(len(members)), nil
}

func (l *local) SetCard(ctx context.Context, key string) (int64, error) {
	value, ok := l.data.Load(key)
	if !ok {
		return 0, errors.New("not exist")
	}
	item := value.(*MapItem)
	vs := item.V.([]any)
	return int64(len(vs)), nil
}

func (l *local) SetRem(ctx context.Context, key string, members ...any) (int64, error) {
	value, ok := l.data.Load(key)
	if !ok {
		return 0, errors.New("not exist")
	}
	item := value.(*MapItem)
	vs := item.V.([]any)
	for _, member := range members {
		for i2, v := range vs {
			if reflect.DeepEqual(member, v) {
				vs = append(vs[:i2], vs[i2+1:])
				break
			}
		}
	}
	return int64(len(vs)), nil
}

func (l *local) SetIsMember(ctx context.Context, key string, member any) (bool, error) {
	value, ok := l.data.Load(key)
	if !ok {
		return false, errors.New("not exist")
	}
	item := value.(*MapItem)
	vs := item.V.([]any)
	for _, v := range vs {
		if reflect.DeepEqual(member, v) {
			return true, nil
		}
	}
	return false, nil
}
func (l *local) expiration() {
	for {
		l.data.Range(func(key, value any) bool {
			v, ok := value.(*MapItem)
			if ok {
				if time.Now().Sub(v.CreateAt) >= v.TTL {
					l.data.Delete(key)
				}
			}
			return true
		})
		time.Sleep(1 * time.Second)
	}
}

func NewLocalClient() constrain.IRedis {
	l := &local{data: sync.Map{}}
	go l.expiration()
	return l
}
