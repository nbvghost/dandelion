package config

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/nbvghost/dandelion/library/util"
	"log"
	"os"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/go-redis/redis/v8"
)

type MicroServerConfig struct {
	MicroServer MicroServer
	Port        int
	IP          string
	Addr        string
}

func (m *MicroServerConfig) UnRegister() error {
	fileName, err := m.MicroServer.getFileName()
	if err != nil {
		return err
	}

	err = os.Remove(fileName)
	if err != nil {
		return err
	}
	return nil
}
func (m *MicroServerConfig) Register() error {
	var err error
	if m.IP == "" {
		m.IP = util.NetworkIP()
		if m.IP == "" {
			return errors.New("无法获取本机ip")
		}
	}
	if m.Port == 0 {
		m.Port, err = util.RandomNetworkPort()
		if err != nil {
			return err
		}
	}

	fileName, err := m.MicroServer.getFileName()
	if err != nil {
		return err
	}

	address := fmt.Sprintf("%s:%d", m.IP, m.Port)

	_, err = os.Stat(fileName)
	if err == nil {
		err = os.Remove(fileName)
		if err != nil {
			return err
		}
	}

	err = writeLock(fileName, []byte(address))
	if err != nil {
		return err
	}
	return nil
}

func NewMicroServerConfig(microServerKey MicroServer, port int, ip string) *MicroServerConfig {
	return &MicroServerConfig{
		MicroServer: microServerKey,
		Port:        port,
		IP:          ip,
	}
}

type RedisOptions struct {
	// The network type, either tcp or unix.
	// Default is tcp.
	Network string
	// host:port address.
	Addr string

	// Use the specified Username to authenticate the current connection
	// with one of the connections defined in the ACL list when connecting
	// to a Redis 6.0 instance, or greater, that is using the Redis ACL system.
	Username string
	// Optional password. Must match the password specified in the
	// requirepass server configuration option (if connecting to a Redis 5.0 instance, or lower),
	// or the User Password when connecting to a Redis 6.0 instance, or greater,
	// that is using the Redis ACL system.
	Password string

	// Database to be selected after connecting to the server.
	DB int

	// Maximum number of retries before giving up.
	// Default is 3 retries; -1 (not 0) disables retries.
	MaxRetries int
	// Minimum backoff between each retry.
	// Default is 8 milliseconds; -1 disables backoff.
	MinRetryBackoff time.Duration
	// Maximum backoff between each retry.
	// Default is 512 milliseconds; -1 disables backoff.
	MaxRetryBackoff time.Duration

	// Dial timeout for establishing new connections.
	// Default is 5 seconds.
	DialTimeout time.Duration
	// Timeout for socket reads. If reached, commands will fail
	// with a timeout instead of blocking. Use value -1 for no timeout and 0 for default.
	// Default is 3 seconds.
	ReadTimeout time.Duration
	// Timeout for socket writes. If reached, commands will fail
	// with a timeout instead of blocking.
	// Default is ReadTimeout.
	WriteTimeout time.Duration

	// Type of connection pool.
	// true for FIFO pool, false for LIFO pool.
	// Note that fifo has higher overhead compared to lifo.
	PoolFIFO bool
	// Maximum number of socket connections.
	// Default is 10 connections per every available CPU as reported by runtime.GOMAXPROCS.
	PoolSize int
	// Minimum number of idle connections which is useful when establishing
	// new connection is slow.
	MinIdleConns int
	// Connection age at which client retires (closes) the connection.
	// Default is to not close aged connections.
	MaxConnAge time.Duration
	// Amount of time client waits for connection if all connections
	// are busy before returning an error.
	// Default is ReadTimeout + 1 second.
	PoolTimeout time.Duration
	// Amount of time after which client closes idle connections.
	// Should be less than server's timeout.
	// Default is 5 minutes. -1 disables idle timeout check.
	IdleTimeout time.Duration
	// Frequency of idle checks made by idle connections reaper.
	// Default is 1 minute. -1 disables idle connections reaper,
	// but idle connections are still discarded by the client
	// if IdleTimeout is set.
	IdleCheckFrequency time.Duration

	// Enables read only queries on slave nodes.
	readOnly bool

	// TLS Config to use. When set TLS will be negotiated.
	TLSConfig *tls.Config

	// Limiter interface used to implemented circuit breaker or rate limiter.
	Limiter redis.Limiter
}

func (m RedisOptions) ToOptions() redis.Options {
	return redis.Options{
		Network:            m.Network,
		Addr:               m.Addr,
		Username:           m.Username,
		Password:           m.Password,
		DB:                 m.DB,
		MaxRetries:         m.MaxRetries,
		MinRetryBackoff:    m.MinRetryBackoff,
		MaxRetryBackoff:    m.MaxRetryBackoff,
		DialTimeout:        m.DialTimeout,
		ReadTimeout:        m.ReadTimeout,
		WriteTimeout:       m.WriteTimeout,
		PoolFIFO:           m.PoolFIFO,
		PoolSize:           m.PoolSize,
		MinIdleConns:       m.MinIdleConns,
		MaxConnAge:         m.MaxConnAge,
		PoolTimeout:        m.PoolTimeout,
		IdleTimeout:        m.IdleTimeout,
		IdleCheckFrequency: m.IdleCheckFrequency,
		TLSConfig:          m.TLSConfig,
		Limiter:            m.Limiter,
	}
}

type ServerConfig struct {
	Server MicroServerConfig
	Etcd   clientv3.Config
	Redis  RedisOptions
}

type PostgresqlConfig struct {
	Host     string `json:"Host"`
	User     string `json:"User"`
	Password string `json:"Password"`
	DBName   string `json:"DBName"`
	Port     int    `json:"Port"`
	SSLMode  string `json:"SSLMode"`
	TimeZone string `json:"TimeZone"`
}

func (m *PostgresqlConfig) GetDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s", m.Host, m.User, m.Password, m.DBName, m.Port, m.SSLMode, m.TimeZone)
}

func GetENV(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}
	log.Println(fmt.Sprintf("env %s %s", key, value))
	return value
}
