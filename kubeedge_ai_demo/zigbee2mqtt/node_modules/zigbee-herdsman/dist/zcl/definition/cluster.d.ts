import DataType from './dataType';
import { ParameterDefinition } from './tstype';
interface AttributeDefinition {
    ID: number;
    type: DataType;
    manufacturerCode?: number;
}
interface ClusterDefinition {
    ID: number;
    manufacturerCode?: number;
    attributes: {
        [s: string]: AttributeDefinition;
    };
    commands: {
        [s: string]: CommandDefinition;
    };
    commandsResponse: {
        [s: string]: CommandDefinition;
    };
}
interface CommandDefinition {
    ID: number;
    parameters: ParameterDefinition[];
    response?: number;
}
declare const Cluster: {
    [s: string]: ClusterDefinition;
};
export default Cluster;
//# sourceMappingURL=cluster.d.ts.map