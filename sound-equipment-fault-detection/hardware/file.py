import os, glob

def read_config_file(file_path="/etc/config/normal"):
    try:
        with open(file_path, 'r') as file:
            return file.read().strip()
    except IOError as e:
        print(f"Failed to read config file: {e}")
        return None
def get_files():
    # Get all files in the input directory
    config_value = read_config_file()
    if config_value == "true":
        input_dir = './data/normal'
    elif config_value == "false":
        input_dir = './data/anomaly'
    else:
        input_dir = './data/normal'
    files = glob.glob(os.path.join(input_dir, '*'))
    return files