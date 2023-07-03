### A Example project written in go demonstrating Graph Object Database (GODB)

#### Proposed Features

- Con-current
- Distributed
- Redundancy
- Failure Recovery
- Transaction
- Deadlock Avoidance

![Image](https://github.com/afmjoaa/GODB/assets/32086430/466826bd-713c-45b4-b5f5-1b7fd069eee7)


#### Serialization Deserialization

We need serialization and deserialization in two steps

1. To save the data to database along with schema. -> **[go-avro](https://github.com/linkedin/goavro)**
2. To Communicate via RPC with among clint-server. -> **[gRPC with Protocol buffer](https://grpc.io/docs/languages/go/quickstart/)**

Brain-teaser: How to gather schema in protocol-buffer from user to generate data-object specific methods. (Future scope)



