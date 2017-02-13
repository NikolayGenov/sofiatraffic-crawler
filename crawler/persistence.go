package crawler

import (
	"encoding/json"

	"github.com/garyburd/redigo/redis"
)

func (s *SofiaTrafficCrawler) SaveLines(conn redis.Conn) (reply interface{}, err error) {
	serialized, err := json.Marshal(s.Lines)
	if err != nil {
		return
	}
	reply, err = conn.Do("SET", "SofiaTraffic/lines", serialized)
	if err != nil {
		return
	}
	return
}

func (s *SofiaTrafficCrawler) SaveSchedules(conn redis.Conn) error {
	serialized, err := json.Marshal(s.Schedules)
	if err != nil {
		return err
	}
	_, err = conn.Do("SET", "SofiaTraffic/schedules", serialized)
	if err != nil {
		return err
	}
	return nil
}

func (s *SofiaTrafficCrawler) LoadLines(conn redis.Conn) error {
	serialized, err := redis.Bytes(conn.Do("GET", "SofiaTraffic/lines"))
	if err != nil {
		return err
	}
	return json.Unmarshal(serialized, &s.Lines)
}
func (s *SofiaTrafficCrawler) LoadSchedules(conn redis.Conn) error {
	serialized, err := redis.Bytes(conn.Do("GET", "SofiaTraffic/schedules"))
	if err != nil {
		return err
	}
	return json.Unmarshal(serialized, &s.Schedules)
}
