package main

import (
	"os"
	"reflect"
	"testing"
)

func Test_getNewOwner_positive(t *testing.T) {
	ip := "91.210.181.37"
	result, _ := getOwner(ip)
	if result != "Digmia s.r.o." {
		t.Error("incorrect result: expected 'Digmia s.r.o.', got", result)
	} else {
		t.Logf("[Test_getOwner] %s -> %s\n", ip, result)
	}

}

func Test_getNewOwner_negative(t *testing.T) {
	ip := "192.168.0.1"
	result, _ := getOwner(ip)
	if result != "Unknown" {
		t.Error("incorrect result: expected 'Unknown', got", result)
	} else {
		t.Logf("[Test_getOwner] %s -> %s\n", ip, result)
	}

}

func Test_colorize(t *testing.T) {

	ip := []string{"37.9.169.172"}
	target := "q6x04kc0oqm8v40t93pl34eg77dy1p0dp.ctdl.ml"
	token := " q6x04kc0oqm8v40t93pl34eg77dy1p0dp" // token has a leading space

	testline, err := os.ReadFile("testdata/testline.txt")
	if err != nil {
		t.Error("testline.txt: Unable to read file.")
	}
	line := string(testline)

	result := colorize(line, ip, target, token)
	if result == line {
		t.Error("incorrect result: no transfoormation on input, got", result)
	} else {
		t.Logf("[Test_colorize][MISS] %s \n", result)
	}

	resultCached := colorize(line, ip, target, token)
	if result == line {
		t.Error("incorrect result: no transfoormation on input, got", result)
	} else {
		t.Logf("[Test_colorize][HIT] %s \n", resultCached)
	}

}

func Test_Cache(t *testing.T) {
	var cacheFile = "testdata/cache.json"
	var testCache1 = make(ipcache)
	var testCache2 = make(ipcache)

	// Load cache from test file
	_ = cacheLoad(testCache1, cacheFile)
	// Save it
	cacheSave(testCache1, cacheFile)
	// Load it again
	_ = cacheLoad(testCache2, cacheFile)

	result := reflect.DeepEqual(testCache1, testCache2)

	if !result {
		t.Error("incorrect result: cache data does not match", result)
	} else {
		t.Log("[Test_Cache] OK:")
	}
}
