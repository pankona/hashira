package syncutil

import (
	"log"

	"github.com/pankona/hashira/syncclient"
)

func (c *Client) TestAccessToken(accesstoken string) {
	sc := syncclient.New()
	if err := sc.TestAccessToken(accesstoken); err != nil {
		log.Printf("test access token failed: %v", err)
	}
	log.Printf("access token is valid. hashira-web works!")
}
