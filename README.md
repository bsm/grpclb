# grpclb

[![Build Status](https://travis-ci.org/bsm/grpclb.png?branch=master)](https://travis-ci.org/bsm/grpclb)
[![GoDoc](https://godoc.org/github.com/bsm/grpclb?status.png)](http://godoc.org/github.com/bsm/grpclb)
[![Gem Version](https://badge.fury.io/rb/grpclb.svg)](https://badge.fury.io/rb/grpclb)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

External Load Balancing Service solution for gRPC written in Go. The approach follows the
[proposal](https://github.com/grpc/grpc/blob/master/doc/load-balancing.md) outlined by the
core gRPC team.

grpclb load-balancer provides a neutral API which can be integrated with various service discovery
frameworks. An example service discovery implementation is provided for [Consul](discovery/consul/).

## Usage

### Load Balancer

Please also see the bootstrap for [Consul backed load-balancers](cmd/grpc-lb-consul/main.go)
as a reference for building load balancers. Either use the command directly or build your very own.

### Server

Servers can optionally report load to the Load Balancer. An example:

See [Documentation](https://godoc.org/github.com/bsm/grpclb/load)

### Client

See [Documentation](https://godoc.org/github.com/bsm/grpclb#NewResolver)
