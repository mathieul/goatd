require "ffi-rzmq"

class TeamsController < ApplicationController
  include ActionController::MimeResponds

  respond_to :json

  def index
    list = []
    respond_with list
  end
end
