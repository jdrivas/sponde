package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Connection contains necessary data to connect to the JupytherHub API
// 	HubURL - the connection end point
//  token - the Token needed for Authorization.

type Connection struct {
	Name   string
	HubURL string
	Token  string
}

// For viper maps
const (
	hubURLKey = "huburl"
	tokenKey  = "token"
)

// State referenced below, and only below.
var currentConnection *Connection
var defaultConnection Connection

// SetConnection sets the connection to the named connection.
// Or returns an error if the named connection doesn't exist.
// Connection is not set if there is an error, and whathever
//  connection is current will continued to be used.
func SetConnection(cName string) (err error) {
	if cName == "default" {
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

// GetConnectionName return the connection's name
func GetConnectionName() string {
	return getCurrentConnection().Name
}

// GetHubURL returns the connection HUBUrl
func GetHubURL() string {
	return getCurrentConnection().HubURL
}

func GetToken() string {
	return getCurrentConnection().Token
}

// GetToken returns  the connections Token
func GetSafeToken() string {
	return MakeSafeTokenString(getCurrentConnection())
}

func MakeSafeTokenString(c Connection) string {
	token := "*****"
	if !viper.GetBool("neverShowToken") {
		if viper.GetBool("showToken") {
			token = c.Token
			if token == "" {
				token = "<empty-token>"
			}
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
		if k == "default" {
			noDefault = false
		}
		cons = append(cons, k)
	}
	if noDefault {
		cons = append(cons, "default")
	}
	return cons
}

func GetConnections() []Connection {
	consMap := getConnectionMap()
	cons := []Connection{}
	for k, v := range consMap {
		if k != "default" {
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
		// hubURL := connsMap.(map[string]interface{})[huURLKey]
		// token := connsMap.(map[string]interface{})[tokenKey]
		conn = Connection{name, hubURL, token}
	}
	return conn, ok
}

func getConnectionMap() map[string]interface{} {
	return viper.GetStringMap("connections")
}

// InitConnections sets up the default connection, sets the current Connection to default,
// and should be called whenever the Viper config file gets
// reloaded.
// Provide a defaultHubURL, if no is provided, the http://127.0.0.1:8081 will be used.
func InitConnections(defaultHubURL string) {
	conn, ok := getConnectionByName("default")
	if ok {
		setDefaultConnection(conn)
	} else {

		hubURL := defaultHubURL
		token := ""

		defaultName := viper.GetString("defaultConnection")
		if defaultName != "" {
			dnc, ok := getConnectionByName(defaultName)
			if ok {
				hubURL = dnc.HubURL
				token = dnc.Token
			}
		}

		conn = Connection{"default", hubURL, token}
		setDefaultConnection(conn)
	}
	// Only set if it's not areadly been set.
	// We want it to be durable when in interactive
	if currentConnection == nil {
		setCurrentConnection(conn)
	}
}
