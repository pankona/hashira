package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/pankona/hashira/api"
)

func main() {
	port := os.Getenv("HASHIRA_API_SERVER_PORT")
	if port == "" {
		port = "18082"
		log.Printf("HASHIRA_AUTH_PORT is not specified. Use default port: %s", port)
	}

	taskStore := &taskStore{
		taskMapByUserID:  map[string]map[string]api.Task{},
		priorityByUserID: map[string]api.Priority{},
	}
	api := &api.API{TaskStore: taskStore}

	http.Handle("/api/v1/upload", &upload{api: api})
	http.Handle("/api/v1/download", &download{api: api})

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

type taskStore struct {
	taskMapByUserID  map[string]map[string]api.Task
	priorityByUserID map[string]api.Priority
}

func (s *taskStore) SaveTasks(userID string, ts api.Tasks) error {
	taskMap, ok := s.taskMapByUserID[userID]
	if !ok {
		// create a new task map for specified user
		s.taskMapByUserID[userID] = map[string]api.Task{}
		taskMap = s.taskMapByUserID[userID]
	}

	for _, v := range ts {
		taskMap[v.ID] = api.Task{
			ID:        v.ID,
			Name:      v.Name,
			Place:     v.Place,
			IsDeleted: v.IsDeleted,
		}
	}

	buf, err := json.Marshal(taskMap)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("/tmp/"+userID+"_tasks.json", buf, 0644)
	if err != nil {
		return err
	}

	log.Printf("len of tasks: %v", len(taskMap))
	log.Printf("%v", taskMap)

	return nil
}

func (s *taskStore) SavePriority(userID string, p api.Priority) error {
	priority, ok := s.priorityByUserID[userID]
	if !ok {
		// create a new priority array for specified user
		s.priorityByUserID[userID] = api.Priority{}
		priority = s.priorityByUserID[userID]
	}

	for place, IDs := range p {
		priority[place] = IDs
	}

	buf, err := json.Marshal(priority)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("/tmp/"+userID+"_priority.json", buf, 0644)
	if err != nil {
		return err
	}

	log.Printf("len of priority: %v", len(priority))

	return nil
}

func (s *taskStore) LoadTasks(userID string) (api.Tasks, error) {
	tasks, ok := s.taskMapByUserID[userID]
	if !ok {
		return api.Tasks{}, fmt.Errorf("tasks of [%s] not found", userID)
	}
	return tasks, nil
}

func (s *taskStore) LoadPriority(userID string) (api.Priority, error) {
	priority, ok := s.priorityByUserID[userID]
	if !ok {
		return api.Priority{}, fmt.Errorf("priority of [%s] not found", userID)
	}
	return priority, nil
}

type upload struct {
	api *api.API
}

type Task struct {
	ID        string
	Name      string
	Place     string
	IsDeleted bool
}

type Priority map[string][]string

// priority's key should be one of following strings:
// "BACKLOG", "TODO", "DOING", "DONE"
type UploadRequest struct {
	Tasks    map[string]Task `json:"tasks"`
	Priority Priority        `json:"priority"`
}

func (u *upload) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v", r.URL.Path)

	switch r.Method {
	case http.MethodPost:
		accesstoken, ok := r.Header["Authorization"]
		if !ok || len(accesstoken) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		log.Printf("Authorization: %v", accesstoken)

		me, err := GetMe(accesstoken[0])
		if err != nil {
			log.Printf("failed to get user info from auth service: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		buf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		req := UploadRequest{}
		err = json.Unmarshal(buf, &req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tasks := api.Tasks{}
		for k, v := range req.Tasks {
			tasks[k] = api.Task{
				ID:        v.ID,
				Name:      v.Name,
				Place:     api.Place(v.Place),
				IsDeleted: v.IsDeleted,
			}
		}

		priority := api.Priority{}
		for k, v := range req.Priority {
			priority[api.Place(k)] = v
		}

		err = u.api.Upload(me.ID, tasks, priority)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	default:
		// unsupported
		w.WriteHeader(http.StatusNotImplemented)
	}
}

type download struct {
	api *api.API
}

type DownloadResponse UploadRequest

func (d *download) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		accesstoken, ok := r.Header["Authorization"]
		if !ok || len(accesstoken) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		log.Printf("Authorization: %v", accesstoken)

		me, err := GetMe(accesstoken[0])
		if err != nil {
			log.Printf("failed to get user info from auth service: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tasks, priority, err := d.api.Download(me.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resp := DownloadResponse{
			Tasks: func() map[string]Task {
				ret := map[string]Task{}
				for k, v := range tasks {
					ret[k] = Task{
						ID:        v.ID,
						Name:      v.Name,
						Place:     string(v.Place),
						IsDeleted: v.IsDeleted,
					}
				}
				return ret
			}(),
			Priority: func() Priority {
				ret := Priority{}
				for k, v := range priority {
					ret[string(k)] = v
				}
				return ret
			}(),
		}

		buf, err := json.Marshal(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(buf)
		if err != nil {
			log.Print(err)
		}

	default:
		w.WriteHeader(http.StatusNotImplemented)
	}
}

type User struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	GoogleID       string `json:"google_id"`
	TwitterID      string `json:"twitter_id"`
	AccessToken    string `json:"access_tokens"`
	TwitterIDToken string `json:"twitter_id_token"`
	GoogleIDToken  string `json:"google_id_token"`
}

const authServiceURI = "http://localhost:8080/api/v1"

func GetMe(accesstoken string) (User, error) {
	req, err := http.NewRequest(http.MethodGet, authServiceURI+"/me", nil)
	if err != nil {
		return User{}, err
	}
	req.Header["Authorization"] = []string{accesstoken}

	cli := http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		return User{}, err
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return User{}, err
	}

	log.Printf("%v", string(buf))

	u := User{}
	err = json.Unmarshal(buf, &u)
	if err != nil {
		return User{}, err
	}

	return u, nil
}
