package kafka

import (
	"fmt"
	"github.com/GCFactory/dbo-system/platform/config"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/IBM/sarama"
	"strings"
	"sync"
)

// pool of producers that ensure transactional-id is unique.
type ProducerProvider struct {
	transactionIdGenerator int32

	producersLock sync.Mutex
	producers     []sarama.AsyncProducer

	producerProvider func() sarama.AsyncProducer

	logger logger.Logger
}

func NewKafkaProducer(cfg *config.Config, logger logger.Logger) *ProducerProvider {
	provider := &ProducerProvider{
		logger: logger,
	}
	provider.producerProvider = func() sarama.AsyncProducer {
		saramaCfg := sarama.NewConfig()
		saramaCfg.Version = sarama.DefaultVersion
		saramaCfg.Producer.Idempotent = true
		saramaCfg.Producer.Return.Errors = false
		saramaCfg.Producer.RequiredAcks = sarama.WaitForAll
		saramaCfg.Producer.Partitioner = sarama.NewRoundRobinPartitioner
		saramaCfg.Producer.Transaction.Retry.Backoff = 10
		saramaCfg.Producer.Transaction.ID = "txn_producer_notification"
		saramaCfg.Net.MaxOpenRequests = 1
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
	return provider
}

func (p *ProducerProvider) borrow() (producer sarama.AsyncProducer) {
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

func (p *ProducerProvider) release(producer sarama.AsyncProducer) {
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

func (p *ProducerProvider) clear() {
	p.producersLock.Lock()
	defer p.producersLock.Unlock()

	for _, producer := range p.producers {
		producer.Close()
	}
	p.producers = p.producers[:0]
}

func (p *ProducerProvider) ProduceRecord(topic string, message []byte) error {
	producer := p.borrow()
	defer p.release(producer)

	// Start kafka transaction
	err := producer.BeginTxn()
	if err != nil {
		return err
	}

	// Produce some records in transaction
	var i int64
	// Need concurrency
	for i = 0; i < 1; i++ {
		producer.Input() <- &sarama.ProducerMessage{Topic: topic, Key: nil, Value: sarama.ByteEncoder(message)}
	}

	// commit transaction
	err = producer.CommitTxn()
	if err != nil {
		p.logger.Warnf("Producer: unable to commit txn %s\n", err)
		for {
			if producer.TxnStatus()&sarama.ProducerTxnFlagFatalError != 0 {
				// fatal error. need to recreate producer.
				p.logger.Fatalf("Producer: producer is in a fatal state, need to recreate it")
				break
			}
			// If producer is in abortable state, try to abort current transaction.
			if producer.TxnStatus()&sarama.ProducerTxnFlagAbortableError != 0 {
				err = producer.AbortTxn()
				if err != nil {
					// If an error occured just retry it.
					p.logger.Errorf("Producer: unable to abort transaction: %+v", err)
					continue
				}
				break
			}
			// if not you can retry
			err = producer.CommitTxn()
			if err != nil {
				p.logger.Warnf("Producer: unable to commit txn %s\n", err)
				continue
			}
		}
		return nil
	}
	return nil
}
