# hook

Link tools like Prometheus Alertmanager to Kafka through webhook config.    
Hook listens on localhost to expose an entrypoint to a specific topic via HTTP call.

Hook relies on [kafka-go](https://github.com/segmentio/kafka-go.git) 

## v0.1.3

Prepare config file in the same dir :  
```
server:
  bind: 127.0.0.1:5000
  validate: true

kafka:
  username: toto
  password: pwd
  topic: default_topic
  brokers: localhost:9092
  tls:
    capath: /tmp/capath
```

Run hook service :  
```
hook
```

## TODO

Stability : 
- improve webserver

New features : 
- check JSON format 

Improvements : 
- specify config file instead of local dir
- add SASL SCRAM compatibility
- make SASL_SSL not mandatory

Doc : 
- add examples

Next ?
- retrieve target topic from incoming json
- retrieve SASL user from incoming json

## DONE

Security : 
- TLS => provide a CA path to establish a secure session
- SASL => provide a user/pwd and bind as SASL_PLAINTEXT

Features :
- load config from yaml file 

Stability : 
- enforce sync writes into topic => check write kafka consistency OK
- handle errors => http status code OK 
- log to a file => log to stdout instead in json format OK
- add flag validate to test kafka connection before running webserver OK