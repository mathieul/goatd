class TeamsController < ApplicationController
  include ActionController::MimeResponds

  respond_to :json

  def index
    list = nil
    TcpClient.new(AppConfig[:atd_address]) do
      send_request("team", "list", {})
      list = receive_response
    end
    respond_with list
  end
end
