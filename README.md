# kiko
Kiko general purpose recommendation system

---

[![Go](https://github.com/hvuhsg/kiko/actions/workflows/go.yml/badge.svg)](https://github.com/hvuhsg/kiko/actions/workflows/go.yml)


### Key consepts

**Node:** Entity that can have connections to other entities, for example player in a sport group or color in a clothing shop.

**Connection:** The desired distance between entity A and entity B.

**Optimization:** The system is constently trying to update its nodes placement to be as close as possible to the nodes connections. This process is called optimization.

### Usege example
If we have a clothing shop and want to know what items to offer to a customer by its purches history, We can create node per every product in the store and create nodes for sizes (S, M, L, XL) colors, gender, categories, etc... then connect every node to its apropriate size color category etc... after that we will create node for every customer and connect him to its purcest items nodes.
After all of that process we will have the option of just asking the system for the 5 most recommended nodes for the customer node and offering it to him.

### API
The system support REST API as well as gRPC server.

REST: OpenAPI [schema](/communication/rest/spec.yaml)
gRPC: ProtoBuf [schema](/communication/grpc/spec.proto)
The repo also contains [insomnia](insomnia.rest) workspace with all of the requests and enviroments.