# Guardfile

guard 'shell', :all_on_start => false do
  watch(/(.*?)(_test|).go/) { |m| `go test #{m[1]}_test.go` }
end
