package syncutil

import (
	"fmt"

	"github.com/pankona/hashira/sync"
)

func (c *Client) TestAccessToken(accesstoken string) error {
	sc := sync.NewClient()
	if err := sc.TestAccessToken(accesstoken); err != nil {
		return fmt.Errorf("test access token failed: %w", err)
	}
	return nil
}
