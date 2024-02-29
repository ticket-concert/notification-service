package kafka

import (
	"context"

	k "gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

// Producer is collection of function of kafka producer
type Producer interface {
	Publish(topic string, message []byte, kafkaPartition *int32)

	Close(ctx context.Context) error
}

// Consumer is collection of function of kafka consumer
type Consumer interface {
	SetHandler(handler ConsumerHandler)
	Subscribe(topics ...string)

	Close(ctx context.Context) error
}

// ConsumerHandler is a collection of function for handling kafka message
type ConsumerHandler interface {
	SendEmailOtpRegister(message *k.Message, topic string)
	SendEmailTicket(message *k.Message, topic string)
}

///

type KafkaConfig struct {
	username string
	password string
	address  string
}

var kafkaConfig KafkaConfig

func InitKafkaConfig(kafkaUrl string, username string, password string) {
	kafkaConfig = KafkaConfig{
		address:  kafkaUrl,
		username: username,
		password: password,
	}
}

func GetConfig() KafkaConfig {
	return kafkaConfig
}

func (kc KafkaConfig) GetKafkaConfig(groupId string, commit bool) *k.ConfigMap {
	kafkaCfg := k.ConfigMap{}
	if groupId != "" {
		kafkaCfg.SetKey("enable.auto.commit", commit)
	}
	if kc.username != "" {
		kafkaCfg.SetKey("sasl.mechanism", "SCRAM-SHA-512")
		kafkaCfg.SetKey("sasl.username", kc.username)
		kafkaCfg.SetKey("sasl.password", kc.password)
		kafkaCfg.SetKey("security.protocol", "SASL_PLAINTEXT")
		kafkaCfg.SetKey("enable.ssl.certificate.verification", false)
		//switch securityProtocol {
		//case "sasl_ssl":
		//	kafkaCfg.SetKey("security.protocol", securityProtocol)
		//	kafkaCfg.SetKey("ssl.endpoint.identification.algorithm", "https")
		//	kafkaCfg.SetKey("enable.ssl.certificate.verification", true)
		//	break
		//case "sasl_plaintext":
		//	kafkaCfg.SetKey("security.protocol", securityProtocol)
		//	break
		//default:
		//	kafkaCfg.SetKey("security.protocol", "sasl_plaintext")
		//	break
		//}
	}

	kafkaCfg.SetKey("bootstrap.servers", kc.address)
	kafkaCfg.SetKey("group.id", groupId)
	kafkaCfg.SetKey("retry.backoff.ms", 500)
	kafkaCfg.SetKey("socket.max.fails", 10)
	kafkaCfg.SetKey("reconnect.backoff.ms", 200)
	kafkaCfg.SetKey("reconnect.backoff.max.ms", 5000)
	kafkaCfg.SetKey("request.timeout.ms", 5000)
	kafkaCfg.SetKey("partition.assignment.strategy", "roundrobin")
	kafkaCfg.SetKey("auto.offset.reset", "latest")

	return &kafkaCfg
}
