package jupyterhub

// Connection is the data required to talk to a JuptyterHub hub.
// Connection contains necessary data to connect to the JupytherHub API
// HubURL - the connection end point
// token - the Token needed for Authorization.
// and a name for identification.
type Connection struct {
	Name   string
	HubURL string
	Token  string
	Auth   Auth
}

// Auth holds paramaters to handle OAuth outhorization commands.
type Auth struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

// YAML variables. These show up in viper, but the application
// sets and uses them using the config.
const (
	clientIDKey     = "clientID"
	clientSecretKey = "clientSecret"
	redirectURLKey  = "redirectURL"
)

// UpdateAuth sets only non-blank ("") values of auth.
func (a Auth) UpdateAuth(update Auth) {
	if update.ClientID != "" {
		a.ClientID = update.ClientID
	}
	if update.ClientSecret != "" {
		a.ClientSecret = update.ClientSecret
	}
	if update.RedirectURL != "" {
		a.RedirectURL = update.RedirectURL
	}
}

// InitAuth gets defaul values from  the connection
func InitAuth(connectionName string) {
	panic("InitAuth not implemented yet.")
}
