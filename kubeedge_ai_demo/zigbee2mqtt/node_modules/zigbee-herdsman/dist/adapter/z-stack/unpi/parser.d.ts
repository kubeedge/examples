/// <reference types="node" />
import * as stream from 'stream';
declare class Parser extends stream.Transform {
    private buffer;
    constructor();
    _transform(chunk: Buffer, _: string, cb: () => void): void;
    private parseNext;
}
export default Parser;
//# sourceMappingURL=parser.d.ts.map