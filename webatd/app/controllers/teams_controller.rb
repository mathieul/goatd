require "ffi-rzmq"

class TeamsController < ApplicationController
  include ActionController::MimeResponds

  respond_to :json

  def index
    list = [{name: send_message("Blah")}]
    respond_with list
  end

  private

  def send_message(message)
    logger.debug ">>> #send_message(#{message.inspect}"
    ctx = ZMQ::Context.create(1)
    req = ctx.socket(ZMQ::REQ)
    logger.debug ">>> before connect"
    rc = req_sock.connect('tcp://127.0.0.1:5000')
    raise "Failed to connect REQ socket" unless ZMQ::Util.resultcode_ok?(rc)
    logger.debug ">>> before send_string"
    rc = req_sock.send_string(message)
    raise "Failed to send message" if error_check(rc)
    rep = ''
    logger.debug ">>> before recv_string"
    rc = req_sock.recv_string(rep)
    raise "Failed to receive response" if error_check(rc)
    req_sock.close
    ctx.terminate
    rep
  end
end
