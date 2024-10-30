import time

def read_config_file(file_path="/etc/config/normal"):
    try:
        with open(file_path, 'r') as file:
            return file.read().strip()
    except IOError as e:
        print(f"Failed to read config file: {e}")
        return None