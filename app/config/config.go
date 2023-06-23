package config

import (
	"github.com/abtergo/abtergo/pkg/website"
)

type Config struct {
	Website website.Config `env:"WEBSITE"`
}
