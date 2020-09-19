/// <reference types="node" />
import { Type, Subsystem } from './constants';
declare class Frame {
    readonly type: Type;
    readonly subsystem: Subsystem;
    readonly commandID: number;
    readonly data: Buffer;
    readonly length: number;
    readonly fcs: number;
    constructor(type: Type, subsystem: Subsystem, commandID: number, data: Buffer, length?: number, fcs?: number);
    toBuffer(): Buffer;
    static fromBuffer(length: number, fcsPosition: number, buffer: Buffer): Frame;
    private static calculateChecksum;
    toString(): string;
}
export default Frame;
//# sourceMappingURL=frame.d.ts.map