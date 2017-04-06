import random
import time

import grpc

from . import balancer_pb2, balancer_pb2_grpc


def grpclb_channel(service, grpclb_addr='localhost:8383', credentials=None, options=None):
    """Creates a Channel to a service target.

    Args:
        service: Service name/target.
        grpclb_addr: grpclb server address (host:port).
        credentials: A ChannelCredentials instance to dial the service with.
            The service is connected to insecurely if unset.
        options: A sequence of string-value pairs according to which to configure
            the created channel.

    Returns:
        A Channel to the service through which RPCs may be conducted.
    """
    chan = grpc.insecure_channel(grpclb_addr)
    stub = balancer_pb2_grpc.LoadBalancerStub(chan)
    req = balancer_pb2.ServersRequest(target=service)

    while True:
        resp = stub.Servers(req)
        if len(resp.servers):
            break
        time.sleep(random.random() * 1.5)

    if credentials is None:
        service_chan = grpc.insecure_channel(resp.servers[0].address, options)
    else:
        service_chan = grpc.secure_channel(resp.servers[0].address, credentials, options)

    return service_chan
