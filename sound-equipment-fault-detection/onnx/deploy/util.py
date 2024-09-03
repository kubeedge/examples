import numpy as np

class Wave2Mel(object):
    def __init__(self, sr,
                 n_fft=1024,
                 n_mels=128,
                 win_length=1024,
                 hop_length=512,
                 power=2.0
                 ):
        super(Wave2Mel, self).__init__()
        self.mel_transform = NumpyMelSpectrogram(sample_rate=sr,
                                                                  win_length=win_length,
                                                                  hop_length=hop_length,
                                                                  n_fft=n_fft,
                                                                  n_mels=n_mels,
                                                                  power=power)
        self.amplitude_to_db = NumpyAmplitudeToDB(stype='power')

    def __call__(self, x):
        x = self.mel_transform(x)
        x = self.amplitude_to_db(x)
        return x.astype(np.float32)



class NumpyAmplitudeToDB:
    def __init__(self, stype='power', top_db=80.0):
        """
        Args:
            stype (str): If 'power', assumes input is a power spectrogram (magnitude squared).
                         If 'magnitude', assumes input is a magnitude spectrogram.
            top_db (float): Threshold the output to `top_db` below the max.
        """
        assert stype in ['power', 'magnitude'], "stype must be either 'power' or 'magnitude'"
        self.stype = stype
        self.multiplier = 10.0 if stype == 'power' else 20.0
        self.top_db = top_db

    def __call__(self, spectrogram: np.ndarray) -> np.ndarray:
        """
        Args:
            spectrogram (np.ndarray): Input spectrogram (power or magnitude).

        Returns:
            np.ndarray: Spectrogram in decibels (dB).
        """
        # Convert to decibels (dB)
        amin = 1e-10  # To avoid taking log of zero
        spectrogram_db = self.multiplier * np.log10(np.maximum(amin, spectrogram))

        # If top_db is set, we threshold the values to be no lower than `top_db` below the max
        if self.top_db is not None:
            max_db = spectrogram_db.max()
            spectrogram_db = np.maximum(spectrogram_db, max_db - self.top_db)

        return spectrogram_db

class NumpyMelSpectrogram:
    def __init__(self,
                 sample_rate: int = 16000,
                 n_fft: int = 400,
                 win_length: int = None,
                 hop_length: int = None,
                 f_min: float = 0.,
                 f_max: float = None,
                 pad: int = 0,
                 n_mels: int = 128,
                 power: float = 2.,
                 normalized: bool = False,
                 center: bool = True,
                 pad_mode: str = "reflect",
                 onesided: bool = True,
                 norm: str = None) -> None:
        self.sample_rate = sample_rate
        self.n_fft = n_fft
        self.win_length = win_length if win_length is not None else n_fft
        self.hop_length = hop_length if hop_length is not None else self.win_length // 2
        self.pad = pad
        self.power = power
        self.normalized = normalized
        self.n_mels = n_mels
        self.f_max = f_max if f_max is not None else sample_rate / 2
        self.f_min = f_min
        self.center = center
        self.pad_mode = pad_mode
        self.onesided = onesided
        self.norm = norm

    def _stft(self, waveform: np.ndarray) -> np.ndarray:
        """Short-time Fourier transform (STFT) implementation using numpy."""
        if self.center:
            pad_amount = self.n_fft // 2
            waveform = np.pad(waveform, (pad_amount, pad_amount), mode=self.pad_mode)

        window = np.hanning(self.win_length)
        stft_result = np.array([
            np.fft.rfft(window * waveform[i:i + self.win_length])
            for i in range(0, len(waveform) - self.win_length + 1, self.hop_length)
        ])
        return np.abs(stft_result) ** self.power

    def _mel_filterbank(self) -> np.ndarray:
        """Compute a Mel filterbank."""
        # Frequency bins
        fft_freqs = np.fft.rfftfreq(self.n_fft, 1.0 / self.sample_rate)
        mel_bins = np.linspace(self._hz_to_mel(self.f_min), self._hz_to_mel(self.f_max), self.n_mels + 2)
        hz_bins = self._mel_to_hz(mel_bins)

        # Create filterbank
        mel_filterbank = np.zeros((self.n_mels, len(fft_freqs)))
        for i in range(1, self.n_mels + 1):
            lower = hz_bins[i - 1]
            center = hz_bins[i]
            upper = hz_bins[i + 1]

            left_slope = (fft_freqs >= lower) & (fft_freqs <= center)
            mel_filterbank[i - 1, left_slope] = (fft_freqs[left_slope] - lower) / (center - lower)

            right_slope = (fft_freqs >= center) & (fft_freqs <= upper)
            mel_filterbank[i - 1, right_slope] = (upper - fft_freqs[right_slope]) / (upper - center)

        if self.norm == 'slaney':
            mel_filterbank /= mel_filterbank.sum(axis=1, keepdims=True)

        return mel_filterbank

    def _hz_to_mel(self, hz):
        return 2595 * np.log10(1 + hz / 700)

    def _mel_to_hz(self, mel):
        return 700 * (10 ** (mel / 2595) - 1)

    def __call__(self, waveform: np.ndarray) -> np.ndarray:
        """Compute the Mel Spectrogram."""
        spectrogram = self._stft(waveform)
        mel_filterbank = self._mel_filterbank()
        mel_spectrogram = np.dot(spectrogram, mel_filterbank.T)
        mel_spectrogram = np.transpose(mel_spectrogram, (1, 0))
        return mel_spectrogram