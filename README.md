# hook

Link tools like Prometheus Alertmanager to Kafka through webhook config.    
Hook listens on localhost to expose an entrypoint to a specific topic via HTTP call.

Hook relies on [kafka-go](https://github.com/segmentio/kafka-go.git) 

## v0.1.1

Run hook service :  
```
hook --topic testtopic --brokers localhost:9092 --listen localhost:3000 --validate
```


## TODO

Stability : 
- improve webserver

New features : 
- TLS
- SASL
- check JSON format 

Next ?
- retrieve target topic from incoming json
- retrieve SASL user from incoming json

## DONE

Stability : 
- enforce sync writes into topic => check write kafka consistency OK
- handle errors => http status code OK 
- log to a file => log to stdout instead in json format OK
- add flag validate to test kafka connection before running webserver OK