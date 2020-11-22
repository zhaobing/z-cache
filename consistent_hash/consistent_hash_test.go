package consistent_hash

import (
	"fmt"
	"strconv"
	"testing"
)

func TestHashCircle_AddPhysicalNode(t *testing.T) {
	t.Log("start test")
	var replicas = 3
	hashFn := func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	}

	hashCircle := New(replicas, hashFn)
	hashCircle.AddPhysicalNode("6", "4", "2")
	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}

	for k, v := range testCases {
		if hashCircle.GetPhysicalNode(k) != v {
			t.Errorf("Asking for %s,should have yielded %s", k, v)
		}
	}

	hashCircle.AddPhysicalNode("8")
	testCases["27"] = "8"
	for k, v := range testCases {
		if hashCircle.GetPhysicalNode(k) != v {
			t.Errorf("Asking for %s,should have yielded %s", k, v)
		}
	}

	fmt.Println("start test")

}
