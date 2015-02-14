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
	Status    string    `json:"done"`
}
