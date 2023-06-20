# RedisConn

## Example Redis Pub.Sub

```go
package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/go-redis/redis"
	"github.com/sivaosorg/govm/logger"
	"github.com/sivaosorg/govm/redisx"
	"github.com/sivaosorg/redisconn/redisconn"
)

func main() {
	client, s := redisconn.NewClient(*redisx.GetRedisConfigSample().SetPassword("Tm5@P@ssw0rd").SetDebugMode(false))
	logger.Infof("state connection: %v", s)
	if !s.IsConnected {
		panic(s.Error)
	}
	scanner := bufio.NewScanner(os.Stdin)
	ps := redisconn.NewRedisPubSub(client)
	pubsubs, err := ps.Subscribe("chat")
	if err != nil {
		panic(err)
	}
	go func() {
		for _, pubsub := range pubsubs {
			go func(ps *redis.PubSub) {
				for {
					msg, err := ps.ReceiveMessage()
					if err != nil {
						panic(err)
					}
					logger.Infof("%s: %s", msg.Channel, msg.Payload)
				}
			}(pubsub)
		}
	}()
	logger.Infof("Type your message and hit enter to send a chat message!")
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		if scanner.Text() == "exit" || scanner.Text() == "quit" {
			break
		}
		// Publish the chat message to the "chat" channel
		err := ps.Publish("chat", scanner.Text())
		if err != nil {
			panic(err)
		}
	}
	// Unsubscribe from multiple channels and close the connection
	err = ps.Unsubscribe("chat")
	if err != nil {
		panic(err)
	}
	err = ps.Close()
	if err != nil {
		panic(err)
	}
}
```

```go
package main

import (
	"github.com/go-redis/redis"
	"github.com/sivaosorg/govm/logger"
	"github.com/sivaosorg/govm/redisx"
	"github.com/sivaosorg/redisconn/redisconn"
)

func main() {
	client, s := redisconn.NewClient(*redisx.GetRedisConfigSample().SetPassword("Tm5@P@ssw0rd").SetDebugMode(false))
	logger.Infof("state connection: %v", s)
	if !s.IsConnected {
		panic(s.Error)
	}
	ps := redisconn.NewRedisPubSub(client)
	pubSubs, err := ps.Subscribe("channel_1", "channel_2", "channel_3")
	if err != nil {
		panic(err)
	}
	go func() {
		for _, pubsub := range pubSubs {
			go func(ps *redis.PubSub) {
				for {
					msg, err := ps.ReceiveMessage()
					if err != nil {
						panic(err)
					}
					logger.Infof("%s: %s", msg.Channel, msg.Payload)
				}
			}(pubsub)
		}
	}()

	// Publish a message to the subscribed channels
	err = ps.Publish("channel_1", "content1")
	if err != nil {
		panic(err)
	}

	err = ps.Publish("channel_2", "content2")
	if err != nil {
		panic(err)
	}

	err = ps.Publish("channel_3", "content3")
	if err != nil {
		panic(err)
	}

	err = ps.Publish("channel_1", "content4")
	if err != nil {
		panic(err)
	}

	// Unsubscribe from multiple channels and close the connection
	err = ps.Unsubscribe("channel_1", "channel_2", "channel_3")
	if err != nil {
		panic(err)
	}

	err = ps.Close()
	if err != nil {
		panic(err)
	}
}

```
