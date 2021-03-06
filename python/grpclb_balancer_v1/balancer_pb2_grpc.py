# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
import grpc

from grpclb_balancer_v1 import balancer_pb2 as grpclb__balancer__v1_dot_balancer__pb2


class LoadBalancerStub(object):
  # missing associated documentation comment in .proto file
  pass

  def __init__(self, channel):
    """Constructor.

    Args:
      channel: A grpc.Channel.
    """
    self.Servers = channel.unary_unary(
        '/grpclb.balancer.v1.LoadBalancer/Servers',
        request_serializer=grpclb__balancer__v1_dot_balancer__pb2.ServersRequest.SerializeToString,
        response_deserializer=grpclb__balancer__v1_dot_balancer__pb2.ServersResponse.FromString,
        )


class LoadBalancerServicer(object):
  # missing associated documentation comment in .proto file
  pass

  def Servers(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')


def add_LoadBalancerServicer_to_server(servicer, server):
  rpc_method_handlers = {
      'Servers': grpc.unary_unary_rpc_method_handler(
          servicer.Servers,
          request_deserializer=grpclb__balancer__v1_dot_balancer__pb2.ServersRequest.FromString,
          response_serializer=grpclb__balancer__v1_dot_balancer__pb2.ServersResponse.SerializeToString,
      ),
  }
  generic_handler = grpc.method_handlers_generic_handler(
      'grpclb.balancer.v1.LoadBalancer', rpc_method_handlers)
  server.add_generic_rpc_handlers((generic_handler,))
