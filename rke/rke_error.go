package rke

import (
	"fmt"
)

const rkeErrorTemplate = `
%s

============= RKE outputs ==============

%s
========================================
`

func wrapErrWithRKEOutputs(err error) error {
	if err == nil {
		return nil
	}

	rkeLogLines := rkeLogBuf.String()
	if rkeLogLines == "" {
		return err
	}
	return fmt.Errorf(rkeErrorTemplate, err, rkeLogLines)
}
