package concurrent

import (
	"testing"
	"sync"
	"strconv"
)
var MapObject *ConcurrentMap = NewConcurrentMap(10);

func init(){
	MapObject = NewConcurrentMap(10);
}

type DataValue  struct{
	Id int
	Name string
}

func TestConcurrentMap_AddOrUpdate(b *testing.T) {

	wait := sync.WaitGroup{}
	wait.Add(1);
	go func() {
		defer wait.Done();
		value1 := MapObject.AddOrUpdate("key_1", DataValue{Id: 1, Name:"name_1"}, func(key string) interface{} {
			return DataValue{Id: 0, Name:"name_func_1"}; });

		if (value1 == nil) {
			b.Errorf("%s", "添加value1失败");
		} else {
			value := value1.(DataValue);
			b.Logf("%s", value.Name);
		}

	}();

	wait.Add(1);
	go func() {
		defer wait.Done();
		value2 := MapObject.AddOrUpdate("key_1", DataValue{Id: 2, Name:"name_2"}, func(key string) interface{} {
			return DataValue{Id: 3, Name:"name_func_2"}; });


		if (value2 == nil) {
			b.Errorf("%s", "添加value2失败");
		} else {
			value := value2.(DataValue);
			b.Logf("%s", value.Name);
		}
	}();

	wait.Wait();

}

func TestConcurrentMap_ContainsKey(t *testing.T) {
	isHas := MapObject.ContainsKey("key_1");

	if(isHas == false){
		t.Fatalf("%s","不存在key_1")
	}else{
		t.Logf("%s",strconv.FormatBool(isHas))
	}

	isHas2 := MapObject.ContainsKey("key_2");

	if(isHas2){
		t.Fatalf("%s","存在key_2")
	}else{
		t.Logf("%s",strconv.FormatBool(isHas2))
	}

}
