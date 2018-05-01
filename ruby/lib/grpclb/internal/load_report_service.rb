require 'grpclb_backend_v1/backend_pb'
require 'grpclb_backend_v1/backend_services_pb'

module Grpclb
  module Internal
    class LoadReportService < Grpclb::Backend::V1::LoadReport::Service
      include GRPC::GenericService

      def initialize(pool)
        @pool = pool
      end

      def load(_req, _call)
        Backend::V1::LoadResponse.new(score: @pool.jobs_waiting)
      end
    end
  end
end
