package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "spt"
	app.Usage = "SNS Publish to Topic\n\n A simple CLI that takes input to STDIN and sends publishes it to an SNS Topic"
	app.UsageText = "cat one_pay_load_per_line.txt | spt --topic-arn ..."
	app.Version = "1.0.0"
	app.HideHelp = true
	app.HideVersion = true

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "topic-arn, t",
			Usage: "SNS topic arn",
		},
		cli.StringFlag{
			Name:  "region, r",
			Usage: "Amazon Web Service REGION if different from the standard credential chain or environment",
		},
		cli.BoolFlag{
			Name:  "help, h",
			Usage: "show this help message",
		},
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "log verbosely",
		},
	}

	app.Action = do
	app.Run(os.Args)
}

func do(c *cli.Context) {
	if c.Bool("help") {
		cli.ShowAppHelpAndExit(c, 1)
	}

	if !c.Bool("verbose") {
		log.SetOutput(ioutil.Discard)
	}

	log.Println("creating AWS session")
	var sess *session.Session
	if region := c.String("region"); region != "" {
		sess = session.Must(
			session.NewSession(
				&aws.Config{Region: aws.String(region)},
			),
		)
	} else {
		sess = session.Must(
			session.NewSessionWithOptions(
				session.Options{
					SharedConfigState: session.SharedConfigEnable,
				},
			),
		)
	}

	snsClient := sns.New(sess)
	topicArn := c.String("topic-arn")

	err := publishToTopic(snsClient, os.Stdin, topicArn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "spt encountered an error: %s", err)
		os.Exit(1)
	}

	return
}

func publishToTopic(client *sns.SNS, input io.Reader, arn string) error {
	log.Print("beginning to read from STDIN")
	scanner := bufio.NewScanner(input)

	var count uint64
	for scanner.Scan() {
		str := strings.TrimSpace(scanner.Text())

		if len(str) == 0 {
			continue
		}

		req := &sns.PublishInput{
			TopicArn: aws.String(arn),
			Message:  aws.String(str),
		}

		_, err := client.Publish(req)
		if err != nil {
			return fmt.Errorf("there was an issue publishing to SNS: %w", err)
		}

		count++
	}

	log.Println("published", count, "messages")
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("the scanner encountered an error: %w", err)
	}

	return nil
}
