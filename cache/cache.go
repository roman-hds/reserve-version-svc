package cache

import "sync"

// A thread-safe in-memory key:value cache
type LatestBuildCache struct {
	Builds map[string]string
	sync.Mutex
}

// Add a build version to the cache, replacing any existing build version
func (c *LatestBuildCache) Save(branch, buildVersion string) {
	c.Lock()
	defer c.Unlock()
	c.Builds[branch] = buildVersion
}

// Return a build version for specified branch
func (c *LatestBuildCache) Read(branch string) string {
	c.Lock()
	defer c.Unlock()
	return c.Builds[branch]
}

// Returns true if value for specified branch is present in cache
func (c *LatestBuildCache) HasKey(branch string) bool {
	c.Lock()
	defer c.Unlock()
	_, ok := c.Builds[branch]
	return ok
}
