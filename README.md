# spt: SNS Publish to Topic
spt is a small utility to that publishes data to a SNS topic

## Documentation
```
$ spt --help
NAME:
   spt - SNS Publish to Topic

 A simple CLI that takes input to STDIN and sends publishes it to an SNS Topic

USAGE:
   cat one_pay_load_per_line.txt | spt --topic-arn ...

GLOBAL OPTIONS:
   --topic-arn value, -t value  SNS topic arn
   --region value, -r value     Amazon Web Service REGION if different from the standard credential chain or environment
   --help, -h                   show this help message
   --verbose                    log verbosely
```
