# Guardfile

guard 'shell', :all_on_start => false do
    watch(/(.*?)(_test|).go/) do |m|
    file = "#{m[1]}_test.go"
    if File.exists?(file)
      puts "\n\n>>> RUNNING #{file} >>>\n\n"
      `go test #{file}`
    end
  end
end
