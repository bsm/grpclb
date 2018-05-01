require 'grpclb/internal/load_report_service'

class Grpclb::Server < GRPC::RpcServer
  DUMMY_POOL = Struct.new(:jobs_waiting).new(0)

  def initialize(*args)
    super
    handle(Grpclb::Internal::LoadReportService.new(@pool || DUMMY_POOL))
  end
end
