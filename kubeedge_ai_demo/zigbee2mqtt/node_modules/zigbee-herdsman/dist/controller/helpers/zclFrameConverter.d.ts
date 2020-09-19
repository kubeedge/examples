import { ZclFrame } from '../../zcl';
interface KeyValue {
    [s: string]: number | string;
}
declare function attributeKeyValue(frame: ZclFrame): KeyValue;
declare function attributeList(frame: ZclFrame): Array<string | number>;
export { attributeKeyValue, attributeList, };
//# sourceMappingURL=zclFrameConverter.d.ts.map