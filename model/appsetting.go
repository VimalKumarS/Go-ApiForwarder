package appsetting

//Appsetting configuration type
type AppSetting struct {
	GatewaySetting struct {
		URL      string   `json:"url"`
		Prefixes []string `json:"prefixes"`
		Methods  []string `json:"methods"`
		Headers  struct {
			ClientKey string `json:"ClientKey"`
			AuthToken string `json:"AuthToken"`
		} `json:"headers"`
	} `json:"GatewaySetting"`
	ConnectionStrings []struct {
		Name             string `json:"name"`
		ConnectionString string `json:"connectionString"`
		ProviderName     string `json:"providerName"`
	} `json:"connectionStrings"`
}
