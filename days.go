package days

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	Done      bool      `json:"done"`
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

func decodeTask(r io.ReadCloser) (*Task, error) {
	defer r.Close()
	var task Task
	err := json.NewDecoder(r).Decode(&task)
	return &task, err
}

func listTasks(c appengine.Context) ([]Task, error) {
	tasks := []Task{}
	keys, err := datastore.NewQuery("Task").Ancestor(tasklistkey(c)).Order("Scheduled").GetAll(c, &tasks)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(tasks); i++ {
		tasks[i].ID = keys[i].IntID()
	}
	return tasks, nil
}

func (t *Task) delete(c appengine.Context) error {
	return datastore.Delete(c, t.key(c))
}
func init() {
	http.HandleFunc("/tasks", handler)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "home")
}

func handler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	val, err := handleTasks(c, r)
	if err == nil {
		json.NewEncoder(w).Encode(val)
	}
	if err != nil {
		c.Errorf("Task error : %#v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleTasks(c appengine.Context, r *http.Request) (interface{}, error) {
	switch r.Method {
	case "POST":
		task, err := decodeTask(r.Body)
		if err != nil {
			return nil, err
		}
		return task.save(c)
	case "GET":
		return listTasks(c)
	case "DELETE":
		task, err := decodeTask(r.Body)
		if err != nil {
			return nil, err
		}
		return nil, task.delete(c)
	}
	return nil, fmt.Errorf("method not implemented")
}
