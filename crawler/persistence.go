package crawler

import (
	"encoding/json"

	"github.com/garyburd/redigo/redis"
)

//saveLines serializes the list of all lines as json and then sets to a key SofiaTraffic/lines in redis
func (s *SofiaTrafficCrawler) saveLines() {
	conn := s.redisPool.Get()
	defer conn.Close()
	serialized, _ := json.Marshal(s.Lines)
	conn.Do("SET", "SofiaTraffic/lines", serialized)
}

//saveSchedules serializes schedules map as json and then sets to a key SofiaTraffic/schedules in redis
func (s *SofiaTrafficCrawler) saveSchedules() {
	conn := s.redisPool.Get()
	defer conn.Close()
	serialized, _ := json.Marshal(s.Schedules)
	conn.Do("SET", "SofiaTraffic/schedules", serialized)

}

//saveVirtualTableStops serializes the list of virtual table stops as json
// and then sets to a key SofiaTraffic/vtstops in redis
func (s *SofiaTrafficCrawler) saveVirtualTableStops() {
	conn := s.redisPool.Get()
	defer conn.Close()
	serialized, _ := json.Marshal(s.VirtualTableStops)
	conn.Do("SET", "SofiaTraffic/vtstops", serialized)
}

//loadLines de-serializes the list of all lines from json back to list and loads it into
// SofiaTrafficCrawler.Lines field
func (s *SofiaTrafficCrawler) loadLines() {
	conn := s.redisPool.Get()
	defer conn.Close()
	serialized, _ := redis.Bytes(conn.Do("GET", "SofiaTraffic/lines"))
	json.Unmarshal(serialized, &s.Lines)
}

//loadSchedules de-serializes the map of all schedules from json back to map and loads it into
// SofiaTrafficCrawler.Schedules field
func (s *SofiaTrafficCrawler) loadSchedules() {
	conn := s.redisPool.Get()
	defer conn.Close()
	serialized, _ := redis.Bytes(conn.Do("GET", "SofiaTraffic/schedules"))
	json.Unmarshal(serialized, &s.Schedules)
}

//loadVirtualTableStops de-serializes the list of all virtual table stops from json back to a list
// and loads it into SofiaTrafficCrawler.VirtualTableStops field
func (s *SofiaTrafficCrawler) loadVirtualTableStops() {
	conn := s.redisPool.Get()
	defer conn.Close()
	serialized, _ := redis.Bytes(conn.Do("GET", "SofiaTraffic/vtstops"))
	json.Unmarshal(serialized, &s.VirtualTableStops)
}
