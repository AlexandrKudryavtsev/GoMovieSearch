package elastic

import "time"

type Option func(*Elastic)

func Addresses(addrs []string) Option {
	return func(c *Elastic) {
		c.addresses = addrs
	}
}

func Username(username string) Option {
	return func(c *Elastic) {
		c.username = username
	}
}

func Password(password string) Option {
	return func(c *Elastic) {
		c.password = password
	}
}

func ConnAttempts(attempts int) Option {
	return func(c *Elastic) {
		c.connAttempts = attempts
	}
}

func ConnTimeout(timeout time.Duration) Option {
	return func(c *Elastic) {
		c.connTimeout = timeout
	}
}
