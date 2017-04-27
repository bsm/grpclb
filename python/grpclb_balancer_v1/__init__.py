import grpc

from . import balancer_pb2, balancer_pb2_grpc


class ServiceUnavailableError(Exception):
    """Raised when grpclb has no services registered for a target."""


class Client:
    """Simple wrapper handling service discovery from grpclb."""

    def __init__(self, lb_addr, target, service_stub, lb_creds=None, service_creds=None):
        if lb_creds is None:
            chan = grpc.insecure_channel(lb_addr)
        else:
            chan = grpc.secure_channel(lb_addr, lb_creds)
        self.lb_stub = balancer_pb2_grpc.LoadBalancerStub(chan)

        self.target = target
        self.service_stub = service_stub
        self.service_creds = service_creds

        self.delegated = None

    def __getattr__(self, name):
        return getattr(self.delegated, name)

    def reconnect(self):
        """Calls grpclb for a target service address before connecting to it.

        Raises:
            ServiceUnavailableError: When no servers have been registered for the target.
        """
        req = balancer_pb2.ServersRequest(target=self.target)
        resp = self.lb_stub.Servers(req)
        if not resp.servers:
            raise ServiceUnavailableError('No servers available for target {}'.format(self.target))

        if self.service_creds is None:
            chan = grpc.insecure_channel(resp.servers[0].address)
        else:
            chan = grpc.secure_channel(resp.servers[0].address, self.service_creds)

        self.delegated = self.service_stub(chan)
