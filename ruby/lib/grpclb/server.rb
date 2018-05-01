require 'grpclb/internal/load_meter'
require 'grpclb/internal/load_report_service'

class Grpclb::Server < GRPC::RpcServer
  def initialize(*args)
    @meter = Grpclb::Internal::LoadMeter.new
    super
    handle(Grpclb::Internal::LoadReportService.new(@meter)) # TODO: handle method is overridden, so it'll measure LoadReport calls as well :|
  end

  def handle(service)
    super(wrap_service(service))
  end

  private

  def wrap_service(service)
    service = service.is_a?(Class) ? service.new : service

    service.class.rpc_descs.each_key do |meth|
      meth = GRPC::GenericService.underscore(meth.to_s).to_sym
      impl = service.method(meth)

      service.define_singleton_method(meth) do |*a, &b|
        with_increment { impl.call(*a, &b) }
      end
    end

    service
  end

  def with_increment
    yield
  ensure
    @meter.increment
  end
end
