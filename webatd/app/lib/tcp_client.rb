require "ffi-rzmq"

TcpClientError = Class.new(StandardError)

class TcpClient
  def initialize(address)
    @address = address
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
  client = TcpClient.new("tcp://127.0.0.1:4242")
  client.run do
    send_message(ARGV.first || "Allo la terre")
    received = receive_message
    puts "received: #{received.inspect}"
  end
end
