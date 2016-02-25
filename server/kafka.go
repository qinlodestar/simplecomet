package main

import (
	log "code.google.com/p/log4go"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/wvanbergen/kafka/consumergroup"
	"time"
)

const (
	KafkaPushsTopic                    = "KafkaNewMsgsTopic"
	KAFKA_GROUP_NAME                   = "kafka_topic_push_group"
	OFFSETS_PROCESSING_TIMEOUT_SECONDS = 10 * time.Second
	OFFSETS_COMMIT_INTERVAL            = 10 * time.Second
)

var (
	producer sarama.SyncProducer
)

func InitKafka(kafkaAddrs []string) (err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	producer, err = sarama.NewSyncProducer(kafkaAddrs, config)
	return
}

func pushKafka(userId int64) (err error) {
	user := fmt.Sprintf("{\"userId\":%d}", userId)
	message := &sarama.ProducerMessage{Topic: KafkaPushsTopic, Key: sarama.StringEncoder("newMsgs"), Value: sarama.ByteEncoder([]byte(user))}
	if _, _, err = producer.SendMessage(message); err != nil {
		return
	}
	log.Debug("kafka userId=%d", userId)
	return
}

func popKafka() error {
	log.Debug("init popkafka")
	config := consumergroup.NewConfig()
	config.Offsets.Initial = sarama.OffsetNewest
	config.Offsets.ProcessingTimeout = OFFSETS_PROCESSING_TIMEOUT_SECONDS
	config.Offsets.CommitInterval = OFFSETS_COMMIT_INTERVAL
	config.Zookeeper.Chroot = ""
	kafkaTopics := []string{KafkaPushsTopic}
	zooks := []string{"127.0.0.1:2181"}
	cg, err := consumergroup.JoinConsumerGroup(KAFKA_GROUP_NAME, kafkaTopics, zooks, config)
	if err != nil {
		return err
	}
	go func() {
		for err := range cg.Errors() {
			log.Error("consumer error(%v)", err)
		}
	}()
	go func() {
		for msg := range cg.Messages() {
			log.Info("deal with userId:%s, partitionId:%d, Offset:%d, Key:%s msg:%s", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
		}
	}()
	return nil
}

//func parse(op string, msg []byte) {
//
//}
