package days

import (
	"appengine/aetest"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strings"
	"testing"
	"time"
)

func TestTask(t *testing.T) {
	task1 := Task{ID: 100, Summary: "task1", Content: "content1"}
	task2 := Task{ID: 101, Summary: "task2", Content: "content2"}
	assert := assert.New(t)
	assert.Equal(task1.ID, 100)
	assert.Equal(task2.ID, 101)
	assert.Equal(task1.Summary, "task1")
	assert.Equal(task2.Content, "content2")
}

func genTask(day, summary string) (*Task, error) {
	layout := "02/01/2006"
	t, err := time.Parse(layout, day)
	if err != nil {
		return nil, err
	}
	return &Task{Summary: summary, Scheduled: t}, nil
}
func TestTaskListKey(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	k1 := tasklistkey(c)
	k2 := tasklistkey(c)

	if !k1.Equal(k2) {
		t.Error("Expected keys not to be equal.")
	}
	if k1.Kind() != "Task" {
		t.Error("Expected <Task>, got: ", k1.Kind())
	}
	if k1.Incomplete() {
		t.Error("Expected key not to be incomplete.")
	}
	if k2.IntID() != 0 {
		t.Error("Expected intid to be 0, got: ", k2.IntID())
	}
}
func TestKey(t *testing.T) {
	assert := assert.New(t)
	c, err := aetest.NewContext(nil)
	defer c.Close()
	if err != nil {
		t.Fatal(err)
	}
	t1 := Task{ID: 12345, Summary: "some summary"}
	t2 := Task{ID: 222}
	if _, err := t1.save(c); err != nil {
		t.Fatal(err)
	}
	k1 := t1.key(c)
	assert.Equal(k1.IntID(), 12345)
	if _, err := t2.save(c); err != nil {
		t.Fatal(err)
	}
	k2 := t2.key(c)
	assert.Equal(k2.IntID(), 222)
}
func TestDecodeTask(t *testing.T) {
	var testjson1 = `{"summary": "summary1", "content": "content1"}`
	var testjson2 = `{"summary" : "summary2", "content": "content2"}`
	t1, err := decodeTask(ioutil.NopCloser(strings.NewReader(testjson1)))
	if err != nil {
		t.Fatal(err)
	}
	t2, err := decodeTask(ioutil.NopCloser(strings.NewReader(testjson2)))
	assert := assert.New(t)
	assert.Equal(t1.Summary, "summary1")
	assert.Equal(t1.Content, "content1")
	assert.Equal(t2.Summary, "summary2")
	assert.Equal(t2.Content, "content2")
}
