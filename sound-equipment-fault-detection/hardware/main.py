from pymodbus.server.sync import StartTcpServer
from pymodbus.device import ModbusDeviceIdentification
from pymodbus.datastore import ModbusSlaveContext, ModbusServerContext
from pymodbus.datastore.store import ModbusSequentialDataBlock
from pymodbus.transaction import ModbusRtuFramer
import time
import threading
from util import load_audio_file, create_data_block, prepare_modbus_data_blocks
from file import get_files

# Main parameters
TOTAL_SIZE = 320044 # 320044 bytes
CHUNK_MAX = 65535 # Maximum data block size of Modbus protocol
CHUNK_SIZE = 60000 # Size of each data block
ISREAD = 1 # Enter the value of the register, indicating that it has been read
NOREAD = 0 # Enter the value of the register, indicating that it has not been read
assert CHUNK_SIZE < CHUNK_MAX # Ensure that the data block size is less than the maximum value

# Update Modbus data storage
def update_store(new_data, chunk_id, context):
    # Get the ModbusSlaveContext in the context
    slave_context = context[0]
    store = slave_context.store
    store_ir = store['i']
    store_hr = store['h']
    store_ir.setValues(1, new_data) # Set the holding register
    store_hr.setValues(chunk_id+1, [NOREAD]) # Set the input register to unread
    print(f"Modbus data store updated for chunk_id {chunk_id}.")


# Threads for monitoring and updating data
def chunk_thread(context, chunks, files):
    chunk_id = 0 #  
    file_id = 0
    while True:
        co_values = context[0].getValues(3, chunk_id, count=1)[0] # Read holding registers
        if co_values == ISREAD and chunk_id < len(chunks) - 1:
            chunk_id += 1
            new_data = chunks[chunk_id]
            update_store(new_data, chunk_id, context)

        elif co_values == ISREAD and chunk_id >= len(chunks) - 1:
            file_id = (file_id+1) % len(files)
            chunk_id = 0
            chunks, hr_values = prepare_file(files[file_id])
            new_data = chunks[chunk_id]

            slave_context = context[0]
            store = slave_context.store
            store_ir = store['i']
            store_hr = store['h']
            store_ir.setValues(1, new_data)  # Set the holding register
            store_hr.setValues(1, hr_values)  # Set input register to read

        else:
            time.sleep(0.01)

def prepare_file(filename):
    audio_data = load_audio_file(filename)
    data_blocks = create_data_block(audio_data, 2)  # 2 bytes per block
    modbus_data = prepare_modbus_data_blocks(data_blocks)

    chunks = []
    for i in range(0, len(modbus_data), CHUNK_SIZE):
        chunk = modbus_data[i:i + CHUNK_SIZE]
        chunks.append(chunk)

    hr_values = [ISREAD] * len(chunks)
    hr_values[0] = NOREAD
    return chunks, hr_values
    
def init(filename):
    chunks, hr_values = prepare_file(filename)
    # Create a data store
    store = ModbusSlaveContext(
        ir=ModbusSequentialDataBlock(1, chunks[0]),  # Continuous register
        hr=ModbusSequentialDataBlock(1, hr_values)  # Input register for notification
    )
    context = ModbusServerContext(slaves=store, single=True)

    # Set device information
    identity = ModbusDeviceIdentification()
    identity.VendorName = 'PyModbus'
    identity.ProductCode = 'PM'
    identity.VendorUrl = 'http://github.com/bashwork/pymodbus/'
    identity.ProductName = 'PyModbus Server'
    identity.ModelName = 'PyModbus Server'
    identity.MajorMinorRevision = '1.0'
    return context, identity, chunks

def main():
    files = get_files()
    for f in files:
        print(f"Use file [{f}] to emulating audio devices")
    context, identity, chunks = init(filename = files[0])
   # Start the Modbus server thread
    server_thread = threading.Thread(target=lambda: StartTcpServer(
            context, 
            identity=identity,
            address=("0.0.0.0", 5020)
        )
    )
    server_thread.daemon = True
    server_thread.start()

    # Start the data update thread
    update_thread = threading.Thread(target=lambda:chunk_thread(
            context,
            chunks,
            files,                        
        )
    )
    update_thread.daemon = True
    update_thread.start()

   # Main Thread
    while True:
        print("Modbus server running...")
        time.sleep(1)

if __name__ == "__main__":
    main()