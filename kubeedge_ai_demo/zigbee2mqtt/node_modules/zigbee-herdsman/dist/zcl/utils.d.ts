import { DataType } from './definition';
import * as TsType from './tstype';
declare function IsDataTypeAnalogOrDiscrete(dataType: DataType): 'ANALOG' | 'DISCRETE';
declare function getCluster(key: string | number, manufacturerCode?: number): TsType.Cluster;
declare function getGlobalCommand(key: number | string): TsType.Command;
export { getCluster, getGlobalCommand, IsDataTypeAnalogOrDiscrete, };
//# sourceMappingURL=utils.d.ts.map