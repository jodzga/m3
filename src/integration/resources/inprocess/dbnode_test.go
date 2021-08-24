// +build integration_v2
// Copyright (c) 2021  Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package inprocess

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDBNodeFromConfigFile(t *testing.T) {
	fileName, cleanup := makeTestConfig(t, defaultConfig)
	defer cleanup()

	dbnode, err := NewDBNodeFromConfigFile(fileName, DBNodeOptions{})
	require.NoError(t, err)

	require.NoError(t, dbnode.Close())
}

func TestNewDBNodeFromYAML(t *testing.T) {
	dbnode, err := NewDBNodeFromYAML(defaultConfig, DBNodeOptions{})
	require.NoError(t, err)

	require.NoError(t, dbnode.Close())
}

func TestWaitForBootstrap(t *testing.T) {
	dbnode, err := NewDBNodeFromYAML(defaultConfig, DBNodeOptions{})
	require.NoError(t, err)

	require.NoError(t, dbnode.WaitForBootstrap())

	res, err := dbnode.Health()
	require.NoError(t, err)

	require.Equal(t, true, res.Bootstrapped)

	require.NoError(t, dbnode.Close())
}

const defaultConfig = `
db: {}
`

const metricsConfig = `
db: 
  metrics:
    prometheus:
      handlerPath: /metrics
      listenAddress: 0.0.0.0:0
    sanitization: prometheus
    samplingRate: 1.0
    extended: detailed
`

func makeTestConfig(t *testing.T, config string) (string, func()) {
	fd, err := ioutil.TempFile("", "config.yaml")
	require.NoError(t, err)
	cleanup := func() {
		assert.NoError(t, fd.Close())
		assert.NoError(t, os.Remove(fd.Name()))
	}
	_, err = fd.Write([]byte(config))
	require.NoError(t, err)

	return fd.Name(), cleanup
}
