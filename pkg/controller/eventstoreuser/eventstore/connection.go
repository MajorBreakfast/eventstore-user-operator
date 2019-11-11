package eventstore

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/scylladb/go-set/strset"
)

var requestTimeout = 5 * time.Second

// ConnectionOptions contains the information for connecting to an Event Store
type ConnectionOptions struct {
	URL           string
	AdminUsername string
	AdminPassword string
}

// Connection is used to talk to an Event Store
type Connection struct {
	client *resty.Client
}

// NewConnection creates a new Event Store connection
func NewConnection(options ConnectionOptions) Connection {
	return Connection{
		client: resty.New().
			SetBasicAuth(options.AdminUsername, options.AdminPassword).
			SetHostURL(options.URL),
		// SetTimeout(requestTimeout), // ToDo: Reenable!!
	}
}

// User contains basic info about an Event Store user
type User struct {
	FullName  string   `json:"fullName"`
	LoginName string   `json:"loginName"`
	Groups    []string `json:"groups"`
}

// UserWithPassword contains basic info about an Event Store user
type UserWithPassword struct {
	User
	Password string `json:"password"`
}

// IsEqual compares two users for equality
func (user User) IsEqual(otherUser User) bool {
	return user.FullName == otherUser.FullName &&
		user.LoginName == otherUser.LoginName &&
		strset.New(user.Groups...).IsEqual(strset.New(otherUser.Groups...))
}

// GetUser returns the user as defined on the Event Store
func (conn *Connection) GetUser(loginName string) (User, error) {
	var resBody struct{ Data User }
	res, err := conn.client.R().
		SetResult(&resBody).
		Get("/users/" + loginName)

	if err != nil {
		return User{}, err
	}
	if res.StatusCode() == 404 {
		return User{}, ErrUserNotFound
	}
	if res.StatusCode() != 200 {
		return User{}, fmt.Errorf("Getting user failed due to unexpected status code %v", res.StatusCode())
	}

	return resBody.Data, nil
}

// CreateUser creates a user on the Event Store
func (conn *Connection) CreateUser(user UserWithPassword) error {
	res, err := conn.client.R().
		SetBody(user).
		Post("/users/")

	if err != nil {
		return err
	}
	if res.StatusCode() != 201 {
		return fmt.Errorf("User creation failed due to unexpected status code %v", res.StatusCode())
	}

	return nil
}

// UpdateUser updates a user on the Event Store
func (conn *Connection) UpdateUser(user User) error {
	res, err := conn.client.R().
		SetBody(user).
		Put("/users/" + user.LoginName)

	if err != nil {
		return err
	}
	if res.StatusCode() != 200 {
		return fmt.Errorf("User update failed due to unexpected status code %v", res.StatusCode())
	}

	return nil
}

// ResetPassword resets the password for a user on the Event Store
func (conn *Connection) ResetPassword(username string, password string) error {
	type reqBody struct {
		NewPassword string `json:"newPassword"`
	}

	res, err := conn.client.R().
		SetBody(reqBody{NewPassword: password}).
		Post("/users/" + username + "/command/reset-password")

	if err != nil {
		return err
	}
	if res.StatusCode() != 200 {
		return fmt.Errorf("Password reset failed due to unexpected status code %v", res.StatusCode())
	}

	return nil
}

// DeleteUser deletes the user from the Event Store
func (conn *Connection) DeleteUser(username string) error {
	type reqBody struct{ newPassword string }

	res, err := conn.client.R().
		Delete("/users/" + username)

	if err != nil {
		return err
	}
	if !(res.StatusCode() == 200 || res.StatusCode() == 204 || res.StatusCode() == 404) {
		return fmt.Errorf("User deletion failed due to unexpected status code %v", res.StatusCode())
	}

	return nil
}
