package user

// User is entity to represent a user
type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	GoogleID  string `json:"google_id"`
	TwitterID string `json:"twitter_id"`
}
