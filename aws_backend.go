package soko

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type AWSBackend struct {
	SectionConfig

	client *ec2.EC2
}

func NewAWSBackend(config SectionConfig) (*AWSBackend, error) {
	conf := &aws.Config{Region: aws.String(config["region"])}
	cli := ec2.New(conf)

	return &AWSBackend{
		SectionConfig: config,
		client:        cli,
	}, nil
}

func (b *AWSBackend) Get(serverID string, key string) (string, error) {
	params := &ec2.DescribeTagsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("resource-id"),
				Values: []*string{aws.String(serverID)},
			},
			{
				Name:   aws.String("key"),
				Values: []*string{aws.String(key)},
			},
		},
	}

	tags, err := b.client.DescribeTags(params)
	if err != nil {
		return "", err
	}

	if len(tags.Tags) != 1 {
		return "", fmt.Errorf("Invalid size of key %s: %d tags exist", key, len(tags.Tags))
	}

	return *tags.Tags[0].Value, nil
}

func (b *AWSBackend) Put(serverID string, key string, value string) error {
	params := &ec2.CreateTagsInput{
		Resources: []*string{
			aws.String(serverID),
		},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String(key),
				Value: aws.String(value),
			},
		},
	}
	_, err := b.client.CreateTags(params)
	return err
}

func (b *AWSBackend) Delete(serverID string, key string) error {
	params := &ec2.DeleteTagsInput{
		Resources: []*string{
			aws.String(serverID),
		},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String(key),
				Value: aws.String(""),
			},
		},
	}
	_, err := b.client.DeleteTags(params)
	return err
}