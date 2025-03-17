package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

// QueueClient represents an SQS client
type QueueClient struct {
	client *sqs.Client
}

// NewQueueClient creates a new SQS client
func NewQueueClient(ctx context.Context) (*QueueClient, error) {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %v", err)
	}

	// Create SQS client
	client := sqs.NewFromConfig(cfg)
	return &QueueClient{client: client}, nil
}

// PushMessage sends a message to an SQS queue
func (q *QueueClient) PushMessage(ctx context.Context, queueURL string, message interface{}) error {
	// Convert message to JSON
	messageBody, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	// Send message to SQS
	_, err = q.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(string(messageBody)),
	})
	if err != nil {
		return fmt.Errorf("failed to send message to SQS: %v", err)
	}

	return nil
}

// AckMessage acknowledges a message by deleting it from the queue
func (q *QueueClient) AckMessage(ctx context.Context, queueURL string, receiptHandle string) error {
	fmt.Printf("Acknowledging message: %s on queue: %s\n", receiptHandle, queueURL)
	_, err := q.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueURL),
		ReceiptHandle: aws.String(receiptHandle),
	})
	if err != nil {
		return fmt.Errorf("failed to delete message from SQS: %v", err)
	}

	return nil
}

// PushBatchMessages sends multiple messages to an SQS queue in a batch
func (q *QueueClient) PushBatchMessages(ctx context.Context, queueURL string, messages []interface{}) error {
	// SQS allows a maximum of 10 messages per batch
	const maxBatchSize = 10

	// Process messages in batches
	for i := 0; i < len(messages); i += maxBatchSize {
		end := i + maxBatchSize
		if end > len(messages) {
			end = len(messages)
		}

		batch := messages[i:end]
		entries := make([]types.SendMessageBatchRequestEntry, len(batch))

		// Prepare batch entries
		for j, msg := range batch {
			messageBody, err := json.Marshal(msg)
			if err != nil {
				return fmt.Errorf("failed to marshal message: %v", err)
			}

			// Create unique ID for each message
			id := fmt.Sprintf("msg-%d-%d", i, j)
			entries[j] = types.SendMessageBatchRequestEntry{
				Id:          aws.String(id),
				MessageBody: aws.String(string(messageBody)),
			}
		}

		// Send batch to SQS
		_, err := q.client.SendMessageBatch(ctx, &sqs.SendMessageBatchInput{
			QueueUrl: aws.String(queueURL),
			Entries:  entries,
		})
		if err != nil {
			return fmt.Errorf("failed to send batch messages to SQS: %v", err)
		}
	}

	return nil
}
