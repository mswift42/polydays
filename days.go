package days

import (
	"appengine"
	"appengine/datastore"
	// "encoding/json"
	// "fmt"
	// "io"
	// "net/http"
	"time"
)

// Task - struct for one Task, with a taskID, assigned by the datastore
// a task summary, task content, task status (TODO or DONE), and a scheduled
// date.
type Task struct {
	ID        int64     `json:"id" datastore:"-"`
	Summary   string    `json:"summary"`
	Content   string    `json:"content"`
	Scheduled time.Time `json:"scheduled"`
	Status    string    `json:"done"`
}

func tasklistkey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Task", "default_tasklist", 0, nil)
}

// key - for a given Task t, return the corresponding key to the intId from
// the datastore. If no Id has been stored, let the datastore crteate a
// new key for the task.
func (t *Task) key(c appengine.Context) *datastore.Key {
	if t.ID == 0 {
		return datastore.NewIncompleteKey(c, "Task", tasklistkey(c))
	}
	return datastore.NewKey(c, "Task", "", t.ID, tasklistkey(c))
}

func (t *Task) save(c appengine.Context) (*Task, error) {
	k, err := datastore.Put(c, t.key(c), t)
	if err != nil {
		return nil, err
	}
	t.ID = k.IntID()
	return t, nil
}
