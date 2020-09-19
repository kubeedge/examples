import { Subsystem, Type } from '../unpi/constants';
import { Frame as UnpiFrame } from '../unpi';
import { ZpiObjectPayload } from './tstype';
declare class ZpiObject {
    readonly subsystem: Subsystem;
    readonly command: string;
    readonly commandID: number;
    readonly payload: ZpiObjectPayload;
    readonly type: Type;
    private readonly parameters;
    private constructor();
    static createRequest(subsystem: Subsystem, command: string, payload: ZpiObjectPayload): ZpiObject;
    toUnpiFrame(): UnpiFrame;
    static fromUnpiFrame(frame: UnpiFrame): ZpiObject;
    private static readParameters;
    private createPayloadBuffer;
    isResetCommand(): boolean;
}
export default ZpiObject;
//# sourceMappingURL=zpiObject.d.ts.map