package load

import (
	"github.com/neee333ko/component-base/pkg/json"
	"github.com/neee333ko/log"
	"github.com/redis/go-redis/v9"
)

type Command string

type Notification struct {
	Command Command `json:"command"`
	Payload string  `json:"payload"`
}

const (
	Channel              = "iam.cluster.notifications"
	SecretChangedCommand = "SecretChanged"
	PolicyChangedCommand = "PolicyChanged"
)

func redisSignalHandler(f func(), handle func(Command)) func(i interface{}) {
	return func(i interface{}) {
		msg, _ := i.(*redis.Message)

		var noti *Notification = new(Notification)

		if err := json.Unmarshal([]byte(msg.Payload), noti); err != nil {
			log.Error("unmarshal redis pubsub message failed.")
			return
		}

		switch noti.Command {
		case SecretChangedCommand, PolicyChangedCommand:
			log.Info("update signal received.")
			reloadChan <- f
		default:
			log.Warn("unrecognized command.")
		}

		if handle != nil {
			handle(noti.Command)
		}
	}
}
