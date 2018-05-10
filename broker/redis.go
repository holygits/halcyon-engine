package broker

import (
  "time"

  "github.com/gomodule/redigo/redis"
)

// Redis provides the broker.Broker interface over a redis server
type Redis struct {
  *redis.Pool
}

// NewRedis returns a new Redis instance
func NewRedis(addr string) (*Redis, error) {
  p := &redis.Pool{
    MaxIdle: 3,
    IdleTimeout: 120 * time.Second,
    Dial: func () (redis.Conn, error) { 
      return redis.Dial("tcp", addr,
                        redis.DialReadTimeout(10*time.Second),
                        redis.DialWriteTimeout(10*time.Second))
    },
  }

  return &Redis{
    p,
  }, nil

}

// Publish pushes msg onto q
func (r *Redis) Publish(q string, msg []byte) error {
  conn := r.Get()
  defer conn.Close()

  conn.Do("PUBLISH", q, msg)
  conn.Flush()

  return nil
}

// Subscribe starts a listener on q. New messages are returned as a broker.Message over channel
func (r *Redis) Subscribe(q string) (<-chan *Message, error) {
  conn := r.Get()
  defer conn.Close()

  // Subscribe with redis server
  conn.Send("SUBSCRIBE", q)
  conn.Flush()

  var notify := make(chan []byte)
  // Listen for new messages
  go func() {
    for {

    msg, err := c.Receive()
    if err != nil {
      notify <- &Message{Error: err,}
      continue
    }

    buf, err := redis.Bytes(msg)
    if err != nil {
      notify <- &Message{Error: err,}
      continue
    }

    notify <- &Message{Body: buf}
    }
  }

  return notify, nil
}
