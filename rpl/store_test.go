package rpl

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestGoLevelDBStore(t *testing.T) {
	// Create a test dir
	dir, err := ioutil.TempDir("", "wal")
	if err != nil {
		t.Fatalf("err: %v ", err)
	}
	defer os.RemoveAll(dir)

	// New level
	l, err := NewGoLevelDBStore(dir, 0)
	if err != nil {
		t.Fatalf("err: %v ", err)
	}
	defer l.Close()

	testLogs(t, l)
}

func testLogs(t *testing.T, l LogStore) {
	// Should be no first index
	idx, err := l.FirstID()
	if err != nil {
		t.Fatalf("err: %v ", err)
	}
	if idx != 0 {
		t.Fatalf("bad idx: %d", idx)
	}

	// Should be no last index
	idx, err = l.LastID()
	if err != nil {
		t.Fatalf("err: %v ", err)
	}
	if idx != 0 {
		t.Fatalf("bad idx: %d", idx)
	}

	// Try a filed fetch
	var out Log
	if err := l.GetLog(10, &out); err.Error() != "log not found" {
		t.Fatalf("err: %v ", err)
	}

	// Write out a log
	log := Log{
		ID:   1,
		Data: []byte("first"),
	}
	for i := 1; i <= 10; i++ {
		log.ID = uint64(i)
		if err := l.StoreLog(&log); err != nil {
			t.Fatalf("err: %v", err)
		}
	}

	// Attempt to write multiple logs
	for i := 11; i <= 20; i++ {
		nl := &Log{
			ID:   uint64(i),
			Data: []byte("first"),
		}

		if err := l.StoreLog(nl); err != nil {
			t.Fatalf("err: %v", err)
		}
	}

	// Try to fetch
	if err := l.GetLog(10, &out); err != nil {
		t.Fatalf("err: %v ", err)
	}

	// Try to fetch
	if err := l.GetLog(20, &out); err != nil {
		t.Fatalf("err: %v ", err)
	}

	// Check the lowest index
	idx, err = l.FirstID()
	if err != nil {
		t.Fatalf("err: %v ", err)
	}
	if idx != 1 {
		t.Fatalf("bad idx: %d", idx)
	}

	// Check the highest index
	idx, err = l.LastID()
	if err != nil {
		t.Fatalf("err: %v ", err)
	}
	if idx != 20 {
		t.Fatalf("bad idx: %d", idx)
	}
}
