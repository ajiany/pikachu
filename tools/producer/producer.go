package producer

import (
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

type IProducer interface {
	SendMessage(topic string, value []byte) error
}

type Producer struct {
	syncProducer sarama.AsyncProducer
}

type producerConfig struct {
	*sarama.Config
}

func defaultConfig() *producerConfig {
	cfg := new(producerConfig)
	cfg.Config = sarama.NewConfig()

	return cfg
}

type optionFn func(cfg *producerConfig)

func NewProducer(hosts []string, opts ...optionFn) *Producer {
	cfg := defaultConfig()

	for _, opt := range opts {
		opt(cfg)
	}

	sp, err := sarama.NewAsyncProducer(hosts, cfg.Config)
	if err != nil {
		panic(err)
	}
	producer := &Producer{syncProducer: sp}
	go func() {
		for err := range sp.Errors() {
			logrus.WithError(err).Error("producer raise error")
		}
	}()
	return producer
}

func (p *Producer) Close() error {
	return p.syncProducer.Close()
}

func (p *Producer) SendMessage(topic string, value []byte) error {
	select {
	case p.syncProducer.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(value),
	}:
	}
	return nil
}
