# kiko
Kiko general purpose recommendation system

---

[![Go](https://github.com/hvuhsg/kiko/actions/workflows/go.yml/badge.svg)](https://github.com/hvuhsg/kiko/actions/workflows/go.yml)


### Key concepts
**Node**: Entity that can have connections to other entities, for example, a player in a sports group or color in a clothing shop.  
**Connection**: The desired distance between entity A and entity B.  
**Optimization**: The system is constantly trying to update its node placement to be as close as possible to the node's connections. This process is called optimization.  

### Usage example
If we have a clothing shop and want to know what items to offer to a customer by its purchase history, We can create a node for every product in the store and create nodes for sizes (S, M, L, XL) colors, gender, categories, etc... then connect every node to its appropriate size color category, etc... after that we will create a node for every customer and connect it to the customer purchased items nodes. After all of that process, we will have the option to just ask the system for the 5 most recommended nodes for the customer node and offer them to the customer.  

### API
The system support REST API as well as a gRPC server.  
REST: OpenAPI [schema](/communication/rest/spec.yaml).  
gRPC: ProtoBuf [schema](/communication/grpc/spec.proto).  
The repo also contains an [insomnia workspace](insomnia.rest) (a better looking postman) with all of the requests and enviroments.  