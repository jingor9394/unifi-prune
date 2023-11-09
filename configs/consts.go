package configs

const (
	ModelUDMPro     = "UDMPro"
	ModelUDR        = "UDR"
	ModelController = "Controller"
)

const (
	ControllerLoginPath         = "api/login"
	ControllerLogoutPath        = "api/logout"
	ControllerClientActivePath  = "v2/api/site/default/clients/active"
	ControllerClientHistoryPath = "v2/api/site/default/clients/history"
	ControllerCmdRemovalPath    = "api/s/default/cmd/stamgr"
	ControllerWlanConfigPath    = "/v2/api/site/default/wlan/enriched-configuration"
)

const (
	apiPrefix         = "proxy/network/"
	LoginPath         = "api/auth/login"
	LogoutPath        = "api/auth/logout"
	ClientActivePath  = apiPrefix + ControllerClientActivePath
	ClientHistoryPath = apiPrefix + ControllerClientHistoryPath
	CmdRemovalPath    = apiPrefix + ControllerCmdRemovalPath
	WlanConfigPath    = apiPrefix + ControllerWlanConfigPath
)
