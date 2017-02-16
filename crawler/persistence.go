package crawler

import (
	"encoding/json"
	"time"

	"github.com/garyburd/redigo/redis"
)

func newPool(address string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 360 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", address) }}
}

func (s *SofiaTrafficCrawler) saveLines() {
	conn := s.redisPool.Get()
	defer conn.Close()
	serialized, _ := json.Marshal(s.Lines)
	conn.Do("SET", "SofiaTraffic/lines", serialized)
}

func (s *SofiaTrafficCrawler) saveSchedules() {
	conn := s.redisPool.Get()
	defer conn.Close()
	serialized, _ := json.Marshal(s.Schedules)
	conn.Do("SET", "SofiaTraffic/schedules", serialized)

}

func (s *SofiaTrafficCrawler) saveVirtualTableStops() {
	conn := s.redisPool.Get()
	defer conn.Close()
	serialized, _ := json.Marshal(s.VirtualTableStops)
	conn.Do("SET", "SofiaTraffic/vtstops", serialized)
}

func (s *SofiaTrafficCrawler) loadLines() {
	conn := s.redisPool.Get()
	defer conn.Close()
	serialized, _ := redis.Bytes(conn.Do("GET", "SofiaTraffic/lines"))
	json.Unmarshal(serialized, &s.Lines)

}

func (s *SofiaTrafficCrawler) loadSchedules() {
	conn := s.redisPool.Get()
	defer conn.Close()
	serialized, _ := redis.Bytes(conn.Do("GET", "SofiaTraffic/schedules"))
	json.Unmarshal(serialized, &s.Schedules)
}

func (s *SofiaTrafficCrawler) loadVirtualTableStops() {
	conn := s.redisPool.Get()
	defer conn.Close()
	serialized, _ := redis.Bytes(conn.Do("GET", "SofiaTraffic/vtstops"))
	json.Unmarshal(serialized, &s.VirtualTableStops)
}
