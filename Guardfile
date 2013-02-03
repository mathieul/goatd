# Guardfile

guard 'compass', :compile_on_start => true do
  watch(%r{^webapp/app/sass/(.*)\.s[ac]ss})
end

`compass compile`

guard 'coffeescript', :input => 'webapp/app/coffeescripts',
                      :output => 'webapp/tmp/compiled',
                      :all_on_start => true

%w[e2e unit].each do |dir|
  guard 'coffeescript', :input => "webapp/test/#{dir}",
                        :output => "webapp/test/javascripts/#{dir}",
                        :all_on_start => true
end


guard :jammit, :config_path => 'assets.yml' do
  watch(%r{^webapp/tmp/compiled/(.*)\.js$})
  watch('assets.yml')
end

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
