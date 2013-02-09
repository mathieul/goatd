require "ffi-rzmq"

class TeamsController < ApplicationController
  include ActionController::MimeResponds

  respond_to :json

  def index
    message = send_message(params[:msg] || "Allo la terre?")
    list = [{name: message}]
    respond_with list
  end

  private

  def send_message(message)
    ctx = ZMQ::Context.create(1)
    req = ctx.socket(ZMQ::REQ)
    rc = req.connect('tcp://127.0.0.1:5000')
    raise "Failed to connect REQ socket" unless ZMQ::Util.resultcode_ok?(rc)
    rc = req.send_string(message)
    raise "Failed to send message" unless ZMQ::Util.resultcode_ok?(rc)
    rep = ''
    rc = req.recv_string(rep)
    raise "Failed to receive response" unless ZMQ::Util.resultcode_ok?(rc)
    req.close
    ctx.terminate
    rep
  end
end
