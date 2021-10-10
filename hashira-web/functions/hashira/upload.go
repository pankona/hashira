package hashira

import (
	"encoding/json"
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

	var tp TaskAndPriority
	err = json.NewDecoder(r.Body).Decode(&tp)
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
		for i, id := range p {
			if ret[id].Place != k {
				priorities[k] = append(p[:i], p[i+1:]...)
			}
		}
	}

	return TaskAndPriority{
		Tasks:    ret,
		Priority: priorities,
	}, nil
}

func mergePriorities(newPriorities, oldPriorities map[string][]string) map[string][]string {
	ret := map[string][]string{"BACKLOG": {}, "TODO": {}, "DOING": {}, "DONE": {}}
	for k := range ret {
		ret[k] = append(newPriorities[k], oldPriorities[k]...)
		ret[k] = unique(ret[k])
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
