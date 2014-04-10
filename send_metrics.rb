require 'statsd'

# 1.times do
#   fork do
    statsd2 = Statsd.new('127.0.0.1', 8000)
    statsd = Statsd.new('127.0.0.1', 8125)

    loop do
      # statsd.gauge('incremental_gauge', "+2")
      # statsd2.gauge('incremental_gauge', "+2")
      statsd.increment('counter_js')
      statsd2.increment('counter_go')
      statsd.gauge('static_gauge_js', "100")
      statsd2.gauge('static_gauge_go', "100")
      r = rand(100)
      statsd.timing('timer', r)
      statsd2.timing('timer', r)
      sleep(0.001)
    end
  # end
# end
