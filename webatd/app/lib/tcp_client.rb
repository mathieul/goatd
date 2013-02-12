require "ffi-rzmq"
require "yajl"
require "pry"

TcpClientError = Class.new(StandardError)

class TcpClient
  def initialize(address, &block)
    @address = address
    run(&block) if block_given?
  end

  def run(&block)
    connect
    instance_eval(&block)
    req_socket.close
  ensure
    ctx.terminate
  end

  def send_request(service, action, request)
    json = Yajl::Encoder.encode(request)
    req_socket.send_string(json, ZMQ::SNDMORE)
    req_socket.send_string(action, ZMQ::SNDMORE)
    rc = req_socket.send_string(service)
    raise TcpClientError.new("Can't send request") unless ok?(rc)
  end

  def receive_response
    message = ""
    rc = req_socket.recv_string(message)
    raise TcpClientError.new("Can't receive message") unless ok?(rc)
    message
  end

  private

  def ok?(result)
    ZMQ::Util.resultcode_ok?(result)
  end

  def ctx
    @ctx ||= ZMQ::Context.new
  end

  def req_socket
    @req_socket ||= ctx.socket(ZMQ::REQ)
  end

  def connect
    rc = req_socket.connect(@address)
    raise TcpClientError.new("Can't connect to #{@address}") unless ok?(rc)
    self
  end
end

if $0 == __FILE__
  Thread.abort_on_exception = true

  TcpClient.new("tcp://127.0.0.1:4242") do
    services = %w[overview.list team.list team.create]
    # services = %w[overview.list team.list]
    request = case (service = services.sample)
              when "overview.list", "team.list"
                {}
              when "team.create"
                name = ARGV.first || "Bonjour"
                {name: name}
              end
    args = service.split(".")
    args << request
    send_request(*args)
    received = receive_response
    puts "received from #{service.inspect}: #{received}"
  end
end
