package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Connection contains necessary data to connect to the JupytherHub API
// HubURL - the connection end point
// token - the Token needed for Authorization.
// and a name for identification.
type Connection struct {
	Name   string
	HubURL string
	Token  string
}

// YAML Variables which show up in viper, but controlled here.
const (
	hubURLKey                = "huburl"
	tokenKey                 = "token"
	neverShowTokensKey       = "neverShowTokens"
	defaultConnectionNameKey = "defaultConnection"
)

// These ARE controlled here, but require binding to a viper configureationa varialbles,
// so we share the name for the key here. The actual variable show be obtained
// through config rather than in viper
const (
	showTokensKey = "showTokens"
)

const (
	defaultConnectionName = "default"
)

// State referenced below, and only below.
// Public access is through functions.
var currentConnection *Connection
var defaultConnection Connection
var showTokens *bool
var showTokensOnce bool

// SetConnection sets the connection to the named connection.
// Or returns an error if the named connection doesn't exist.
// Connection is not set if there is an error, and whathever
//  connection is current will continued to be used.
func SetConnection(cName string) (err error) {
	if cName == defaultConnectionName {
		setCurrentConnection(getDefaultConnection())
	} else {
		conn, ok := getConnectionByName(cName)
		if ok {
			setCurrentConnection(conn)
		} else {
			err = fmt.Errorf("couldn't find connection \"%s\"", cName)
		}
	}
	return err
}

// GetConnectionName return the current conenctions connection's name
func GetConnectionName() string {
	return getCurrentConnection().Name
}

// GetHubURL returns the current connections HubURL
func GetHubURL() string {
	return getCurrentConnection().HubURL
}

// GetToken gets the current connections Token
func GetToken() string {
	return getCurrentConnection().Token
}

// GetSafeToken returns the connections Token. Will
// Return an empty string if useEmpty is true
// otherwise returns something to display.
func GetSafeToken(useEmpty, useShowTokensOnce bool) string {
	return MakeSafeTokenString(getCurrentConnection(), useEmpty, useShowTokensOnce)
}

// MakeSafeTokenString returns a display vlaue for the token that
// is controlled by the state variables ShowTokens
func MakeSafeTokenString(c Connection, useEmpty bool, useShowTokensOnce bool) (token string) {
	token = c.Token
	if c.Token == "" {
		token = "<enpty-token>"
	} else {
		token = "****"
		if useEmpty {
			token = ""
		}
		show := GetShowTokens()
		if useShowTokensOnce {
			show = show || getShowTokensOnce()
		}
		if !viper.GetBool(neverShowTokensKey) && show {
			token = c.Token
		}
	}
	return token
}

// UpdateDefaultHubURL sets the hubURL for the default Connection
func UpdateDefaultHubURL(hubURL string) {
	defaultConnection.HubURL = hubURL
}

// UpdateDefaultToken sets the toke for the default Connection
func UpdateDefaultToken(token string) {
	defaultConnection.Token = token
}

// GetConnectionNames returns a list of names of defined connections.
func GetConnectionNames() []string {
	consMap := getConnectionMap()
	cons := []string{}
	noDefault := true
	for k := range consMap {
		if k == defaultConnectionName {
			noDefault = false
		}
		cons = append(cons, k)
	}
	if noDefault {
		cons = append(cons, defaultConnectionName)
	}
	return cons
}

// GetConnections returns the defined named connections to a JupyterHub Hub.
func GetConnections() []Connection {
	consMap := getConnectionMap()
	cons := []Connection{}
	for k, v := range consMap {
		if k != defaultConnectionName {
			hubURL, token := getMapValues(v)
			c := Connection{k, hubURL, token}
			cons = append(cons, c)
		}
	}
	cons = append(cons, defaultConnection)
	return cons
}

// Internal API
func getDefaultConnection() Connection {
	return defaultConnection
}

func setDefaultConnection(conn Connection) {
	defaultConnection = conn
}

// SetCurrentConneciton sets the connection that config will use
// for connection variables.
func setCurrentConnection(conn Connection) {
	currentConnection = &conn
}

// GetCurrentConnection returns the current connection object for
// the JupyterhHub API
func getCurrentConnection() Connection {
	return *currentConnection
}

func getMapValues(cm interface{}) (hubURLString, tokenString string) {
	hubURLString = cm.(map[string]interface{})[hubURLKey].(string)
	tokenString = cm.(map[string]interface{})[tokenKey].(string)
	return hubURLString, tokenString
}

func getConnectionByName(name string) (conn Connection, ok bool) {
	connsMap, ok := getConnectionMap()[name]
	if ok {
		hubURL, token := getMapValues(connsMap)
		conn = Connection{name, hubURL, token}
	}
	return conn, ok
}

func getConnectionMap() map[string]interface{} {
	return viper.GetStringMap("connections")
}

//
// Show Tokens
//

// State variable to show connection tokens on the next call
// this will reset to false once the connection has been displayed.
func ShowTokensOnce() {
	showTokensOnce = true
}

// Reset the ShowTokensOnce state to false.
func ResetShowTokensOnce() {
	showTokensOnce = false
}

func getShowTokensOnce() bool {
	return showTokensOnce
}

// Current values of the ShowTokensState
func GetShowTokens() bool {
	return *showTokens
}

// Set the status of the ShowTokens state.
func SetShowTokens(st bool) {
	*showTokens = st
}

// InitConnections sets up the default connection, sets the current Connection to default,
// initializes the ShowTokens state, and should be called whenever the Viper config file gets reloaded.
// Provide a defaultHubURL, if no is provided, the http://127.0.0.1:8081 will be used.
func InitConnections(defaultHubURL string) {
	conn, ok := getConnectionByName(defaultConnectionName)
	if ok {
		setDefaultConnection(conn)
	} else {

		hubURL := defaultHubURL
		token := ""

		defaultName := viper.GetString(defaultConnectionNameKey)
		if defaultName != "" {
			dnc, ok := getConnectionByName(defaultName)
			if ok {
				hubURL = dnc.HubURL
				token = dnc.Token
			}
		}

		conn = Connection{defaultConnectionName, hubURL, token}
		setDefaultConnection(conn)
	}
	// Variables which we only set if not already set..
	// We want it to be durable when in interactive
	if currentConnection == nil {
		setCurrentConnection(conn)
	}

	if showTokens == nil {
		showTokens = new(bool)
		SetShowTokens(viper.GetBool(showTokensKey))
	}

}
