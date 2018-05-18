# grpclb/ruby

[![Gem Version](https://badge.fury.io/rb/grpclb.svg)](https://badge.fury.io/rb/grpclb)

Ruby implementation of grpclb client and server.



## Client

```ruby
require 'grpclb/client' # or just 'grpclb'

# construct client:
client = Grpclb::Client.new('127.0.0.1:8383', 'service-name', HelloWorld::V1::Stub)

# call methods:
client.say_hello(HelloWorld::HelloRequest.new(...))
```

Creates and maintains a load-balanced connection for given service stub.

Automatically reconnects on UNAVAILABLE errors.



## Server

```ruby
require 'grpclb/server' # or just 'grpclb'

# construct server:
server = Grpclb::Server.new(...) # just a subclass of GRPC::RpcServer, same initialize args

# handle services:
server.handle(HelloWorld::V1::ServiceImpl) # or .handle(HelloWorld::V1::ServiceImpl.new)

# set up and run:
server.add_http2_port '127.0.0.1:8080', :this_port_is_insecure
server.run_till_terminated
```

Subclass of `GRPC::RpcServer`, that takes care of grpclb load reporting.
