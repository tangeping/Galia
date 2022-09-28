// Copyright (c) TFG Co. All Rights Reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package conf

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config is a wrapper around a viper config
type Config struct {
	config *viper.Viper
}

// NewConfig creates a new config with a given viper config if given
func NewConfig(cfgs ...*viper.Viper) *Config {
	var cfg *viper.Viper
	if len(cfgs) > 0 {
		cfg = cfgs[0]
	} else {
		cfg = viper.New()
	}

	cfg.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	cfg.AutomaticEnv()
	c := &Config{config: cfg}
	c.fillDefaultValues()
	return c
}

func (c *Config) fillDefaultValues() {
	etcdSDConfig := NewDefaultEtcdServiceDiscoveryConfig()

	defaultsMap := map[string]interface{}{
		// the max buffer size that nats will accept, if this buffer overflows, messages will begin to be dropped
		"cluster.sd.etcd.dialtimeout":              etcdSDConfig.DialTimeout,
		"cluster.sd.etcd.endpoints":                etcdSDConfig.Endpoints,
		"cluster.sd.etcd.prefix":                   etcdSDConfig.Prefix,
		"cluster.sd.etcd.grantlease.maxretries":    etcdSDConfig.GrantLease.MaxRetries,
		"cluster.sd.etcd.grantlease.retryinterval": etcdSDConfig.GrantLease.RetryInterval,
		"cluster.sd.etcd.grantlease.timeout":       etcdSDConfig.GrantLease.Timeout,
		"cluster.sd.etcd.heartbeat.log":            etcdSDConfig.Heartbeat.Log,
		"cluster.sd.etcd.heartbeat.ttl":            etcdSDConfig.Heartbeat.TTL,
		"cluster.sd.etcd.revoke.timeout":           etcdSDConfig.Revoke.Timeout,
		"cluster.sd.etcd.syncservers.interval":     etcdSDConfig.SyncServers.Interval,
		"cluster.sd.etcd.syncserversparallelism":   etcdSDConfig.SyncServers.Parallelism,
		"cluster.sd.etcd.shutdown.delay":           etcdSDConfig.Shutdown.Delay,
		"cluster.sd.etcd.servertypeblacklist":      etcdSDConfig.ServerTypesBlacklist,
	}

	for param := range defaultsMap {
		if c.config.Get(param) == nil {
			c.config.SetDefault(param, defaultsMap[param])
		}
	}
}

// GetDuration returns a duration from the inner config
func (c *Config) GetDuration(s string) time.Duration {
	return c.config.GetDuration(s)
}

// GetString returns a string from the inner config
func (c *Config) GetString(s string) string {
	return c.config.GetString(s)
}

// GetInt returns an int from the inner config
func (c *Config) GetInt(s string) int {
	return c.config.GetInt(s)
}

// GetBool returns an boolean from the inner config
func (c *Config) GetBool(s string) bool {
	return c.config.GetBool(s)
}

// GetStringSlice returns a string slice from the inner config
func (c *Config) GetStringSlice(s string) []string {
	return c.config.GetStringSlice(s)
}

// Get returns an interface from the inner config
func (c *Config) Get(s string) interface{} {
	return c.config.Get(s)
}

// GetStringMapString returns a string map string from the inner config
func (c *Config) GetStringMapString(s string) map[string]string {
	return c.config.GetStringMapString(s)
}

// UnmarshalKey unmarshals key into v
func (c *Config) UnmarshalKey(s string, v interface{}) error {
	return c.config.UnmarshalKey(s, v)
}

// Unmarshal unmarshals config into v
func (c *Config) Unmarshal(v interface{}) error {
	return c.config.Unmarshal(v)
}
