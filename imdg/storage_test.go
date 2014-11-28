package imdg

import (
	"testing"
	"time"
)

func TestDistributingByNodes(t *testing.T) {
	rds := NewStorage(uint32(3))
	rds.AddNode(uint32(0), ":6379")
	rds.AddNode(uint32(1), ":6378")

	node_id := rds.getNodeId("fergusson")
	if node_id != uint32(1) {
		t.Errorf("Unexpected node id - %d for gooud", node_id, 1)
	}
	node_id = rds.getNodeId("cleverley")
	if node_id != uint32(2){
		t.Errorf("Unexpected node id - %d for gooud", node_id, 2)
	}

	node_id = rds.getNodeId("rojo")
	if node_id != uint32(0){
		t.Errorf("Unexpected node id - %d for gooud", node_id, 2)
	}
	_, err := rds.GetNode("cleverley")
	if err == nil {
		t.Errorf("Unexpected founded node for get result")
	}

	_, err = rds.GetNode("rojo")
	if err != nil {
		t.Errorf("Unexpected founded node for get result")
	}

	_, err = rds.GetNode("fergusson")
	if err != nil {
		t.Errorf("Unexpected founded node for get result")
	}
}

func TestGetUnexpectedKeys(t *testing.T){
	rds := NewStorage(uint32(3))
	rds.AddNode(uint32(0), ":6379")
	_, err := rds.Get("rojo")
	if err == nil {
		t.Errorf("Founded key in node")
	}
}

func TestGetExpectedKeys(t *testing.T){
	rds := NewStorage(uint32(3))
	rds.AddNode(uint32(0), ":6379")
	err := rds.Set("rojo", []byte("{\"positioning\": 0}"), time.Second*1)
	if err != nil {
		t.Fatal(err)
	}
	_, err = rds.Get("rojo")
	if err != nil {
		t.Fatal(err)
	}
}
