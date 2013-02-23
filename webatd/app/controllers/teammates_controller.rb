class TeammatesController < ApplicationController
  def index
    list = nil
    TcpClient.new(AppConfig[:atd_address]) do
      send_request("teammate", "list")
      list = receive_response
    end
    render json: list["teammates"]
  end

  def create
    teammate, attributes = nil, params[:teammate]
    TcpClient.new(AppConfig[:atd_address]) do
      send_request("teammate", "create", attributes.slice(:name, :team_uid))
      teammate = receive_response
    end
    render json: teammate, status: :created
  end

  def update
    uid, attributes = params[:id], params[:teammate]
    TcpClient.new(AppConfig[:atd_address]) do
      send_request("teammate", "update", uid: uid, teammate: attributes.slice(:name))
      receive_response
    end
    render json: "", status: :ok
  end

  def destroy
    uid = params[:id]
    TcpClient.new(AppConfig[:atd_address]) do
      send_request("teammate", "destroy", uid: uid)
      receive_response
    end
    render json: "", status: :ok
  end
end
