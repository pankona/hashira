package hashira

type Hashira struct {
	AccessTokenStore     AccessTokenStore
	TaskAndPriorityStore TaskAndPriorityStore
}

func New(atStore AccessTokenStore, tpStore TaskAndPriorityStore) *Hashira {
	return &Hashira{
		AccessTokenStore:     atStore,
		TaskAndPriorityStore: tpStore,
	}
}
