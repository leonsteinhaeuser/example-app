package client

import "fmt"

type ClientMarshalError struct {
	Message string
	Err     error
}

func (c *ClientMarshalError) Error() string {
	return fmt.Errorf("thread client marshal error [%s]: %w", c.Message, c.Err).Error()
}

type ClientRequestError struct {
	Message string
	Err     error
}

func (c *ClientRequestError) Error() string {
	return fmt.Errorf("thread client request error [%s]: %w", c.Message, c.Err).Error()
}
