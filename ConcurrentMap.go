package concurrent

import (
	"sync"
)

const(
	DefaultCapacity = 31
)
type ConcurrentMap struct {
	items map[string]*KeyValuePair
	rwMutex *sync.RWMutex
}
type KeyValuePair struct  {
	Key string
	Value interface{}
	index int
}
func NewConcurrentMap(capacity int) *ConcurrentMap{
	return &ConcurrentMap{
		items   : make(map[string]*KeyValuePair,capacity),
		rwMutex : new(sync.RWMutex),
	};
}

//获取Map中元素个数
func (cmap *ConcurrentMap) Count() int {
	if(cmap == nil || cmap.items == nil){
		return 0;
	}
	return len(cmap.items)
}
//获取所有键名
func  (cmap *ConcurrentMap) Keys() []string {
	if(cmap == nil || cmap.items == nil || len(cmap.items) <= 0){
		return make([]string,0);
	}
	cmap.rwMutex.RLock();
	defer cmap.rwMutex.RUnlock();

	keys := make([]string,len(cmap.items));

	for key, _ := range cmap.items {
		keys = append(keys,key);
	}

	return keys;
}
//获取所有的值
func (cmap *ConcurrentMap) Values() []interface{}  {
	if(cmap == nil || cmap.items == nil || len(cmap.items) <= 0){
		return make([]interface{},0);
	}
	cmap.rwMutex.RLock();
	defer cmap.rwMutex.RUnlock();

	values := make([]interface{},len(cmap.items));

	for _,v := range cmap.items{
		values = append(values,v.Value);
	}
	return values;
}

func (cmap *ConcurrentMap) AddRange(keyValues []KeyValuePair){
	cmap.rwMutex.Lock();
	defer cmap.rwMutex.Unlock();

	for _, item := range keyValues{
		item.index = len(cmap.items);
		cmap.items[item.Key] = &item;

	}
}

//获取Map所有的键值对象
func (cmap *ConcurrentMap) ToSlice() []KeyValuePair {
	if(cmap == nil || cmap.items == nil || len(cmap.items) <= 0){
		return make([]KeyValuePair,0);
	}

	entites := make([]KeyValuePair,len(cmap.items));

	cmap.rwMutex.RLock();
	defer cmap.rwMutex.RUnlock();

	for _,v := range cmap.items{
		value := *v;
		entites = append(entites,value);
	}
	return entites;
}

//如果该键不存在则将值添加到Map中，否则调用指定的函数更新Map中的值
// @key 键名
// @value 键值
// @updateValueFactory 获取更新值得函数
func (cmap *ConcurrentMap) AddOrUpdate(key string,value interface{},updateValueFactory func(key string) interface{}) interface{} {

	cmap.rwMutex.Lock();
	defer cmap.rwMutex.Unlock();

	if _,ok := cmap.items[key];ok{
		value = updateValueFactory(key);
	}
	cmap.items[key] = &KeyValuePair{
		Key   : key,
		Value : value,
		index : len(cmap.items),
	};

	return value;
}

//获取或添加，如果键值已经存在则直接获取，否则通过执行函数返回值，并添加到Map中
func (cmap *ConcurrentMap) GetOrAdd(key string,valueFactory func(key string) interface{}) interface{}{

	cmap.rwMutex.Lock();
	defer cmap.rwMutex.Unlock();

	if v,ok := cmap.items[key];ok{
		oldValue := v.Value;
		return oldValue;
	}
	value := valueFactory(key);
	cmap.items[key] = &KeyValuePair{
		Key   : key,
		Value : value,
		index : len(cmap.items),
	};

	return value;
}

//如果成功将键值添加到Map则返回true，如果Map中已存在返回false
func (cmap *ConcurrentMap) TryAdd(key string,value interface{}) bool {
	cmap.rwMutex.Lock();
	defer cmap.rwMutex.Unlock();

	if _,ok := cmap.items[key] ; ok {
		return false;
	}
	cmap.items[key] = &KeyValuePair{
		Key   : key,
		Value : value,
		index : len(cmap.items),
	};
	return true;
}

//尝试获取值
func (cmap *ConcurrentMap) TryGetValue(key string) (value interface{},isOk bool){
	cmap.rwMutex.RLock();
	defer cmap.rwMutex.RUnlock();

	if v,ok := cmap.items[key] ; ok {
		value = v.Value;
		isOk = ok;
		return ;
	}
	return nil,false;
}



//使用提供的比较器比较，如果返回true则使用提供的值替换并返回true否则返回false
func (cmap *ConcurrentMap) TryReplace(key string,value interface{},comparer func(value interface{}) bool) bool {
	cmap.rwMutex.Lock();
	defer cmap.rwMutex.Unlock();

	var comparerResult bool = false;
	var oldValue KeyValuePair;

	if oldValue,ok := cmap.items[key] ; ok {
		v	:= oldValue.Value;

		comparerResult = comparer(v);
	}
	if(comparerResult){
		oldValue.Value = value;
		return true;
	}
	return false;
}

//清空Map中的所有值
func (cmap *ConcurrentMap) Clear(){
	cmap.rwMutex.Lock();
	defer cmap.rwMutex.Unlock();

	cmap.items = make( map[string]*KeyValuePair,DefaultCapacity);
}

//判断指定的键是否在Map中
func (cmap *ConcurrentMap) ContainsKey(key string) bool {

	cmap.rwMutex.RLock();
	defer cmap.rwMutex.RUnlock();

	_,ok := cmap.items[key];
	return ok;
}
//删除指定键的值
func (cmap *ConcurrentMap) Remove(key string) {
	cmap.rwMutex.Lock();
	defer cmap.rwMutex.Unlock();

	delete(cmap.items,key);
}

