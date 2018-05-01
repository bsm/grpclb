# TODO: would be nice to have some tests for it

module Grpclb
  module Internal
    class LoadMeter
      def initialize(period_seconds: 60)
        @period = period_seconds

        @mutex = Mutex.new
        @scored_at = Time.now
        @count = 0
      end

      def increment(inc=1)
        @mutex.synchronize do
          @count += inc
        end
      end

      def score
        @mutex.synchronize do
          now = Time.now
          passed = now - @scored_at
          @scored_at = now

          break 0 if passed == 0

          if passed < @period
            @scored_at -= passed
            break @count * @period / passed
          end

          s = @count * @period / passed
          @count = 0
          s
        end.to_i
      end
    end
  end
end
