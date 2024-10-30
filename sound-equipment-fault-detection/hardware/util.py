def load_audio_file(filename):
    with open(filename, 'rb') as f:
        return f.read()

def create_data_block(data, block_size):
    return [data[i:i + block_size] for i in range(0, len(data), block_size)]

def prepare_modbus_data_blocks(data_blocks):
    data = []
    for block in data_blocks:
        value = int.from_bytes(block, byteorder='big')
        data.append(value)
    return data