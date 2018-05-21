require 'grpclb_balancer_v1/balancer_pb'
require 'grpclb_balancer_v1/balancer_services_pb'

class Grpclb::Client
  attr_reader :target

  DEFAULT_MAX_RECONNECTS = 3

  # @param [String] lb_addr the load balancer address
  # @param [String] target the target service name, as registered in the backend
  # @param [Class] service_stub the service stub class
  def initialize(lb_addr, target, service_stub, lb_creds: :this_channel_is_insecure, service_creds: :this_channel_is_insecure, max_reconnects: DEFAULT_MAX_RECONNECTS) # rubocop:disable Metrics/ParameterLists, Metrics/LineLength
    @lb_stub = Grpclb::Balancer::V1::LoadBalancer::Stub.new(lb_addr, lb_creds)
    @target = target
    @service_stub  = service_stub
    @service_creds = service_creds
    @max_reconnects = max_reconnects

    service_stub.parent.const_get(:Service).rpc_descs.each_key do |meth|
      meth = GRPC::GenericService.underscore(meth.to_s).to_sym

      define_singleton_method meth do |*a, &b|
        with_reconnect { client.send(meth, *a, &b) }
      end
    end
  end

  def reconnect!
    req     = Grpclb::Balancer::V1::ServersRequest.new(target: target)
    servers = @lb_stub.servers(req).servers
    raise "No servers available for target '#{target}'" if servers.empty?

    primary = servers.first.address
    @client = @service_stub.new(primary, @service_creds)
  end

  private

  def client
    @client || reconnect!
  end

  def with_reconnect
    retries = 0
    begin
      yield
    rescue GRPC::BadStatus => e
      raise unless e.code == GRPC::Core::StatusCodes::UNAVAILABLE
      raise if retries >= @max_reconnects

      reconnect!
      retries += 1
      retry
    end
  end
end
