import onnxruntime as ort
import librosa
import re
import glob, os
import numpy as np
from util import Wave2Mel
import time
from file import read_config_file
from control import set_vled_0, set_vled_1

# Hyper parameters are consistent with config.yaml
wav2mel = Wave2Mel(sr=16000, power=2.0, n_fft=1024, n_mels=128, win_length=1024, hop_length=512)

meta2label = {'fan-id_00': 0, 'fan-id_02': 1, 'fan-id_04': 2, 'fan-id_06': 3, 'pump-id_00': 4, 'pump-id_02': 5, 'pump-id_04': 6,
     'pump-id_06': 7, 'slider-id_00': 8, 'slider-id_02': 9, 'slider-id_04': 10, 'slider-id_06': 11, 'ToyCar-id_01': 12,
     'ToyCar-id_02': 13, 'ToyCar-id_03': 14, 'ToyCar-id_04': 15, 'ToyConveyor-id_01': 16, 'ToyConveyor-id_02': 17,
     'ToyConveyor-id_03': 18, 'valve-id_00': 19, 'valve-id_02': 20, 'valve-id_04': 21, 'valve-id_06': 22,
     'fan-id_01': 23, 'fan-id_03': 24, 'fan-id_05': 25, 'pump-id_01': 26, 'pump-id_03': 27, 'pump-id_05': 28,
     'slider-id_01': 29, 'slider-id_03': 30, 'slider-id_05': 31, 'ToyCar-id_05': 32, 'ToyCar-id_06': 33,
     'ToyCar-id_07': 34, 'ToyConveyor-id_04': 35, 'ToyConveyor-id_05': 36, 'ToyConveyor-id_06': 37, 'valve-id_01': 38,
     'valve-id_03': 39, 'valve-id_05': 40}

def transform(filename):
    machine = "valve"
    id_str = "id_00"
    label = meta2label[machine+'-'+id_str]
    x, _ = librosa.core.load(filename, sr=16000, mono=True)
    x_wav = x[: 16000 * 10]
    x_mel = wav2mel(x_wav)
    x_wav = np.expand_dims(x_wav, axis=0)
    x_mel = np.expand_dims(x_mel, axis=0)
    return x_wav, x_mel, label

def numpy_log_softmax(predict_ids, axis=1):
    # Calculate log_softmax
    exp_predict_ids = np.exp(predict_ids - np.max(predict_ids, axis=axis, keepdims=True))
    log_softmax = np.log(exp_predict_ids / np.sum(exp_predict_ids, axis=axis, keepdims=True))

    return log_softmax

# Load ONNX model
onnx_path = "model/STgram-MFN.onnx"
ort_session = ort.InferenceSession(onnx_path)

# ScoreThreshold
ScoreThreshold = 0.01

file = "/etc/data/received.wav"

while True:
    # Check if the file exists
    if not os.path.exists(file):
        time.sleep(0.1) # Sleep for 0.1 seconds for modbus
        print(f'File {file} does not exist, waiting for it to be created')
        continue

    # Record the start time
    start_time = time.time()
    # Convenience files and operations
    print('=' * 65)
    print(f'Processing file')
    x_wav, x_mel, x_label = transform(file)
    x_label = np.array([x_label], dtype=np.int64)
    # Use the model for inference
    predict_ids, feature = ort_session.run(
        None,  # Output name, if None, all outputs are returned
        {"x_wav": x_wav, "x_mel": x_mel, "x_label": x_label}  # Enter a name and value
    )
    # Calculate log_softmax and perform mean, compression and conversion operations
    probs = numpy_log_softmax(predict_ids, axis=1).mean(axis=0).squeeze()

    # Get the output results
    output = probs[x_label[0]]
    output = np.abs(output)
    if output > ScoreThreshold:
        y_pre = 'anomaly'
        set_vled_0()
    else:
        y_pre = 'normal'
        set_vled_1()
    print("output of the ONNX model: ", output)
    print("y_pre of the ONNX model: ", y_pre)

    elapsed_time = time.time() - start_time  # Calculate elapsed time
    if elapsed_time < 5:
        time.sleep(5 - elapsed_time)  # Sleep to ensure total time per file is 5 seconds