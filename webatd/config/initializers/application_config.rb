file_path = Rails.root + "config" + "application.yml"
config = YAML.load_file(file_path)[Rails.env]
AppConfig = HashWithIndifferentAccess.new(config)
