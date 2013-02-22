# Guardfile

guard 'shell', :all_on_start => false do
  runner = ->(m) {
    changed = File.basename(m[0])
    target = case m[1]
           when "app"
             "goatd/app/#{File.dirname(m[2])}"
           when "acceptance"
             m[0]
           end
    puts "\n\n>>> [#{changed}] CHANGED >>> RUNNING #{target} >>>\n\n"
    `go test #{target}`
  }
  watch /(app)\/(.*?)(_test|).go/, &runner
  watch /(acceptance)\/(.*?)(_test|).go/, &runner
end
