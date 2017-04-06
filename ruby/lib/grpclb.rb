require 'grpclb_backend_v1/backend_pb'
require 'grpclb_backend_v1/backend_services_pb'
require 'grpclb_balancer_v1/balancer_pb'
require 'grpclb_balancer_v1/balancer_services_pb'

class Grpclb::Client < ::SimpleDelegator
  attr_reader :target

  # @param [String] lb_addr the load balancer address
  # @param [String] target the target service name, as registered in the backend
  # @param [Class] service_stub the service stub class
  def initialize(lb_addr, target, service_stub, lb_creds: :this_channel_is_insecure, service_creds: :this_channel_is_insecure)
    @lb_stub = Grpclb::Balancer::V1::LoadBalancer::Stub.new(lb_addr, lb_creds)
    @target  = target
    @service_stub  = service_stub
    @service_creds = service_creds

    reconnect!
  end

  def reconnect!
    req     = Grpclb::Balancer::V1::ServersRequest.new(target: target)
    servers = @lb_stub.servers(req).servers
    raise "No servers available for target '#{target}'" if servers.empty?

    primary = servers.first.address
    client  = @service_stub.new(primary, @service_creds)
    __setobj__(client)
  end

end
