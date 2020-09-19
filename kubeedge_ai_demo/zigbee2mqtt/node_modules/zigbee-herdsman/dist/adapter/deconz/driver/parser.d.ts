/// <reference types="node" />
import * as stream from 'stream';
declare class Parser extends stream.Transform {
    private buffer;
    private decoder;
    constructor();
    private onMessage;
    private onError;
    _transform(chunk: Buffer, _: string, cb: Function): void;
}
export default Parser;
//# sourceMappingURL=parser.d.ts.map