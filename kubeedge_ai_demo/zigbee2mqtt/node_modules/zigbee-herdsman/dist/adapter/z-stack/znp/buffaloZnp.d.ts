import { Buffalo, TsType } from '../../../buffalo';
import { BuffaloZnpOptions } from './tstype';
declare class BuffaloZnp extends Buffalo {
    private readListRoutingTable;
    private readListBindTable;
    private readListNeighborLqi;
    private readListNetwork;
    private readListAssocDev;
    read(type: string, options: BuffaloZnpOptions): TsType.Value;
}
export default BuffaloZnp;
//# sourceMappingURL=buffaloZnp.d.ts.map