/// <reference types="node" />
import * as stream from 'stream';
import Frame from './frame';
declare class Writer extends stream.Readable {
    writeFrame(frame: Frame): void;
    writeBuffer(buffer: Buffer): void;
    _read(): void;
}
export default Writer;
//# sourceMappingURL=writer.d.ts.map