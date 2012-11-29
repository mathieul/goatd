# Guardfile

guard 'shell', :all_on_start => false do
  watch(/(.*?)(_test|).go/) { |m| f = "#{m[1]}_test.go"; `go test #{f}` if File.exists?(f) }
end
