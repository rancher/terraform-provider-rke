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
}

func (c *Config) initLogger() {
	if c.LogBuffer == nil {
		c.LogBuffer = &bytes.Buffer{}
	}

	if c.Debug {
		log.SetLevel(log.DebugLevel)
	}

	var writer io.Writer = c.LogBuffer
	//writer := io.MultiWriter(os.Stderr, c.LogBuffer)

	log.SetOutput(writer)
}

func (c *Config) saveRKEOutput(err error) error {
	if len(c.LogFile) > 0 {
		f, errFile := os.OpenFile(c.LogFile, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
		if errFile != nil {
			return fmt.Errorf("Opening logfile %s err:%v", c.LogFile, errFile)
		}
		defer f.Close()
		if _, errFile := f.Write(c.LogBuffer.Bytes()); errFile != nil {
			return fmt.Errorf("Writing logfile %s err:%v", c.LogFile, errFile)
		}
	}

	if err != nil {
		return fmt.Errorf(rkeErrorTemplate, c.LogBuffer.String(), err)
	}

	return nil
}
