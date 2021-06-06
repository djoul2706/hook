# hook

Link tools like Prometheus Alertmanager to Kafka through webhook config.    
Hook listens on localhost to expose an entrypoint to a specific topic via HTTP call.

Hook relies on [kafka-go](https://github.com/segmentio/kafka-go.git) 

## v0.1.0

Run hook service :  
```
hook --topic testtopic --brokers localhost:9092 --listen localhost:3000
```


## TODO

Stability : 
- enforce sync writes into topic 
- handle errors
- log to a file 
- improve webserver

New features : 
- TLS
- SASL
- check JSON format 

Next ?
- retrieve target topic from incoming json
- retrieve SASL user from incoming json

