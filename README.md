# circleci-exporter
Prometheus CircleCI Exporter


Exposes number of deploys per day, to a Prometheus compatible endpoint. 
*This is a very basic implementation, designed for a specific purpose. If you wish to extend/fork this repo to be something greater...  I'm more than open to any pull requests / feedback*

## Configuration

This exporter is setup to take the following parameters from environment variables:
* `URL` The CircleCI insights URL for your workflow
* `AUTH_TOKEN` Auth token to use CircleCI API

## Install and deploy

Build a docker image:
```
docker build -t <image-name> .
docker run --restart=always -p 9179:9179 <image-name>
```

## Metrics

Metrics will be made available on port 9179 by default
