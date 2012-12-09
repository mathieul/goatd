# Guardfile

guard 'shell', :all_on_start => false do
  runner = ->(m) {
    file = "#{m[1]}_test.go"
    if File.exists?(file)
      puts "\n\n>>> RUNNING #{file} >>>\n\n"
      `go test #{file}`
    end
  }
  watch /(app\/.*?)(_test|).go/, &runner
  watch /(acceptance\/.*?)(_test|).go/, &runner
end
