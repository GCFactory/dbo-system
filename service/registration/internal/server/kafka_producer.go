package server

import (
	"fmt"
	"github.com/GCFactory/dbo-system/service/registration/config"
	"github.com/IBM/sarama"
	"strings"
	"sync"
)

// pool of producers that ensure transactional-id is unique.
type KafkaProducerProvider struct {
	transactionIdGenerator int32

	producersLock sync.Mutex
	producers     []sarama.AsyncProducer

	producerProvider func() sarama.AsyncProducer
}

func NewKafkaProducer(cfg *config.Config) (*KafkaProducerProvider, error) {
	saramaCfg := sarama.NewConfig()
	saramaCfg.Version = sarama.DefaultVersion
	saramaCfg.Producer.Idempotent = true
	saramaCfg.Producer.Return.Errors = false
	saramaCfg.Producer.RequiredAcks = sarama.WaitForAll
	saramaCfg.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	saramaCfg.Producer.Transaction.Retry.Backoff = 10
	saramaCfg.Producer.Transaction.ID = "txn_producer"
	saramaCfg.Net.MaxOpenRequests = 1

	provider := &KafkaProducerProvider{}
	provider.producerProvider = func() sarama.AsyncProducer {
		suffix := provider.transactionIdGenerator
		// Append transactionIdGenerator to current saramaCfg.Producer.Transaction.ID to ensure transaction-id uniqueness.
		if saramaCfg.Producer.Transaction.ID != "" {
			provider.transactionIdGenerator++
			saramaCfg.Producer.Transaction.ID = saramaCfg.Producer.Transaction.ID + "-" + fmt.Sprint(suffix)
		}
		producer, err := sarama.NewAsyncProducer(strings.Split(cfg.KafkaProducer.Brokers, ";"), saramaCfg)
		if err != nil {
			return nil
		}
		return producer
	}
	return provider, nil
}

func (p *KafkaProducerProvider) borrow() (producer sarama.AsyncProducer) {
	p.producersLock.Lock()
	defer p.producersLock.Unlock()

	if len(p.producers) == 0 {
		for {
			producer = p.producerProvider()
			if producer != nil {
				return
			}
		}
	}

	index := len(p.producers) - 1
	producer = p.producers[index]
	p.producers = p.producers[:index]
	return
}

func (p *KafkaProducerProvider) release(producer sarama.AsyncProducer) {
	p.producersLock.Lock()
	defer p.producersLock.Unlock()

	// If released producer is erroneous close it and don't return it to the producer pool.
	if producer.TxnStatus()&sarama.ProducerTxnFlagInError != 0 {
		// Try to close it
		_ = producer.Close()
		return
	}
	p.producers = append(p.producers, producer)
}

func (p *KafkaProducerProvider) clear() {
	p.producersLock.Lock()
	defer p.producersLock.Unlock()

	for _, producer := range p.producers {
		producer.Close()
	}
	p.producers = p.producers[:0]
}
