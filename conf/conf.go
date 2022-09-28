package conf

import "time"

var (
	LenStackBuf = 4096

	// log
	LogLevel string
	LogPath  string
	LogFlag  int

	// console
	ConsolePort   int
	ConsolePrompt string = "Leaf# "
	ProfilePath   string

	// cluster
	ListenAddr      string
	ConnAddrs       []string
	PendingWriteNum int
)

// EtcdServiceDiscoveryConfig Etcd service discovery config
type EtcdServiceDiscoveryConfig struct {
	Endpoints   []string
	User        string
	Pass        string
	DialTimeout time.Duration
	Prefix      string
	Heartbeat   struct {
		TTL time.Duration
		Log bool
	}
	SyncServers struct {
		Interval    time.Duration
		Parallelism int
	}
	Revoke struct {
		Timeout time.Duration
	}
	GrantLease struct {
		Timeout       time.Duration
		MaxRetries    int
		RetryInterval time.Duration
	}
	Shutdown struct {
		Delay time.Duration
	}
	ServerTypesBlacklist []string
}

// NewDefaultEtcdServiceDiscoveryConfig Etcd service discovery default config
func NewDefaultEtcdServiceDiscoveryConfig() *EtcdServiceDiscoveryConfig {
	return &EtcdServiceDiscoveryConfig{
		Endpoints:   []string{"localhost:2379"},
		User:        "",
		Pass:        "",
		DialTimeout: time.Duration(5 * time.Second),
		Prefix:      "pitaya/",
		Heartbeat: struct {
			TTL time.Duration
			Log bool
		}{
			TTL: time.Duration(60 * time.Second),
			Log: false,
		},
		SyncServers: struct {
			Interval    time.Duration
			Parallelism int
		}{
			Interval:    time.Duration(120 * time.Second),
			Parallelism: 10,
		},
		Revoke: struct {
			Timeout time.Duration
		}{
			Timeout: time.Duration(5 * time.Second),
		},
		GrantLease: struct {
			Timeout       time.Duration
			MaxRetries    int
			RetryInterval time.Duration
		}{
			Timeout:       time.Duration(60 * time.Second),
			MaxRetries:    15,
			RetryInterval: time.Duration(5 * time.Second),
		},
		Shutdown: struct {
			Delay time.Duration
		}{
			Delay: time.Duration(300 * time.Millisecond),
		},
		ServerTypesBlacklist: nil,
	}
}

// NewEtcdServiceDiscoveryConfig Etcd service discovery config with default config paths
func NewEtcdServiceDiscoveryConfig(config *Config) *EtcdServiceDiscoveryConfig {
	conf := NewDefaultEtcdServiceDiscoveryConfig()
	if err := config.UnmarshalKey("pitaya.cluster.sd.etcd", &conf); err != nil {
		panic(err)
	}
	return conf
}
