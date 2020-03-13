package rke

import (
	"bytes"
	"fmt"
	"io"
	"os"

	//rkelog "github.com/rancher/rke/log"
	log "github.com/sirupsen/logrus"
)

const (
	rkeErrorTemplate = `
============= RKE outputs ==============
%s
%s
========================================
`
)

// Config type of RKE Config
type Config struct {
	Debug     bool
	LogBuffer *bytes.Buffer
	LogFile   string
	File      *os.File
}

func (c *Config) initLogger() {
	if c.LogBuffer == nil {
		c.LogBuffer = &bytes.Buffer{}
	}

	if c.Debug {
		log.SetLevel(log.DebugLevel)
	}

	var writer io.Writer = c.LogBuffer
	if len(c.LogFile) > 0 {
		if c.File == nil {
			f, errFile := os.OpenFile(c.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if errFile != nil {
				log.Errorf("Opening logfile %s err:%v", c.LogFile, errFile)
				return
			}
			c.File = f
		}
		writer = io.MultiWriter(c.LogBuffer, c.File)
	}
	log.SetOutput(writer)
}

func (c *Config) saveRKEOutput(err error) error {
	if c.File != nil {
		defer c.File.Close()
		defer c.File.Sync()
	}
	if err != nil {
		return fmt.Errorf(rkeErrorTemplate, c.LogBuffer.String(), err)
	}

	return nil
}
