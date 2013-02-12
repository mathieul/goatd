class HomeController < ApplicationController
  include ActionController::MimeResponds

  respond_to :json

  def index
    render "index"
  end

  def overview
    list = nil
    TcpClient.new(AppConfig[:atd_address]) do
      send_request("overview", "list", {})
      list = receive_response
    end
    respond_with list
  end
end
