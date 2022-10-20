package config

import (
	"github.com/ajiany/pikachu/tools/producer"
	"sync"
)

var defProducer producer.IProducer
var defProducerOnce sync.Once

// Producer 消费者client初始化
func Producer() producer.IProducer {
	if defProducer != nil {
		return defProducer
	}

	defProducerOnce.Do(func() {
		defProducer = producer.NewProducer(Cfg.KafkaHost)
	})

	return defProducer
}
