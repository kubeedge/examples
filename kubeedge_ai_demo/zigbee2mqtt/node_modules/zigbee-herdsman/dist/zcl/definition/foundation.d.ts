import { ParameterDefinition } from './tstype';
interface FoundationDefinition {
    ID: number;
    parseStrategy: 'repetitive' | 'flat' | 'oneof';
    parameters: ParameterDefinition[];
    response?: number;
}
declare const Foundation: {
    [s: string]: FoundationDefinition;
};
export default Foundation;
//# sourceMappingURL=foundation.d.ts.map