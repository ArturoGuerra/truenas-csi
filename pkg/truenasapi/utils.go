package truenasapi

import "fmt"

/* Returns a valid api url for the truenas api */
func (c *client) parseurl(suffix string) string {
	return fmt.Sprintf("http://%s/api/v2.0/%s", c.Config.Host, suffix)
}
