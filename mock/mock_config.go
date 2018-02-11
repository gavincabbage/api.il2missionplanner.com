package mock

import (
	"github.com/gavincabbage/api.il2missionplanner.com/config"
)

func Config() *config.Config {
	mockConfig := &config.Config{}
	mockConfig.Env = "test"
	mockConfig.Host = "test.com"
	mockConfig.Port = "4321"
	mockServers := make(map[string]string)
	mockServers["testserver"] = "https://test.server.com/test/json.json"
	mockConfig.Servers = mockServers
	return mockConfig
}
