require "ffi-rzmq"

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

  def send_message(message)
    rc = req_socket.send_string(message)
    raise TcpClientError.new("Can't send message") unless ok?(rc)
    self
  end

  def send_typed_message(message, type)
    req_socket.send_string(type, ZMQ::SNDMORE)
    send_message(message)
  end

  def receive_message
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
  TcpClient.new("tcp://127.0.0.1:4242") do
    message = ARGV.first || "Allo la terre"
    services = %w[overview teams]
    service = services.sample
    send_typed_message(service, message)
    received = receive_message
    puts "received from #{service.inspect}: #{received.inspect}"
  end
end
