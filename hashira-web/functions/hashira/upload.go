package hashira

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func (h *Hashira) Upload(w http.ResponseWriter, r *http.Request) {
	accesstoken, err := h.retrieveAccessTokenFromHeader(r.Context(), r.Header)
	if err != nil {
		log.Printf("failed to retrieve accesstoken from header: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	uid, err := h.AccessTokenStore.FindUidByAccessToken(r.Context(), accesstoken)
	if err != nil {
		log.Printf("could not find a user who has the accesstoken %v: %v", accesstoken, err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("failed to read body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	m := map[string]interface{}{}
	if err := json.Unmarshal(buf, &m); err != nil {
		log.Printf("failed to unmarshal: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, ok := m["data"]
	if ok {
		buf, err = json.Marshal(m["data"])
		if err != nil {
			log.Printf("failed to marshal: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	var tp TaskAndPriority
	err = json.Unmarshal(buf, &tp)
	if err != nil {
		log.Printf("failed to read body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	oldtp, err := h.TaskAndPriorityStore.Load(r.Context(), uid)
	if err != nil {
		log.Printf("failed to load task and priorities: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if tp, err = mergeTaskAndPriorities(tp, oldtp); err != nil {
		log.Printf("failed to merge task and priorities: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for k, task := range tp.Tasks {
		if task.IsDeleted {
			delete(tp.Tasks, k)
		}
	}

	if err := h.TaskAndPriorityStore.Save(r.Context(), uid, tp); err != nil {
		log.Printf("failed to save tasks and priorities: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func mergeTaskAndPriorities(newtp, oldtp TaskAndPriority) (TaskAndPriority, error) {
	ret := map[string]Task{}

	for k, v := range oldtp.Tasks {
		ret[k] = v
	}
	for k, v := range newtp.Tasks {
		ret[k] = v
	}

	priorities := mergePriorities(newtp.Priority, oldtp.Priority)

	// Remove priorities if the place is not matched to task's place
	for k, p := range priorities {
		for i := 0; i < len(p); i++ {
			taskID := p[i]
			if ret[taskID].Place != k {
				p = append(p[:i], p[i+1:]...)
				i -= 1
			}
		}
		priorities[k] = p
	}

	return TaskAndPriority{
		Tasks:    ret,
		Priority: priorities,
	}, nil
}

func mergePriorities(newPriorities, oldPriorities map[string][]string) map[string][]string {
	ret := map[string][]string{"BACKLOG": {}, "TODO": {}, "DOING": {}, "DONE": {}}
	for k := range ret {
		ret[k] = mergeStringSlice(newPriorities[k], oldPriorities[k])
	}
	return ret
}

func unique(ss []string) []string {
	keys := make(map[string]struct{})
	ids := []string{}

	for _, id := range ss {
		if _, ok := keys[id]; !ok {
			keys[id] = struct{}{}
			ids = append(ids, id)
		}
	}
	return ids
}
