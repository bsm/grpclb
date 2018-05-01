require 'grpclb_backend_v1/backend_pb'
require 'grpclb_backend_v1/backend_services_pb'

module Grpclb
  module Internal
    class LoadReportService < Grpclb::Backend::V1::LoadReport::Service
      include GRPC::GenericService

      def initialize(meter)
        @meter = meter
      end

      def load(_, _)
        Backend::V1::LoadResponse.new(score: @meter.score)
      end
    end
  end
end
