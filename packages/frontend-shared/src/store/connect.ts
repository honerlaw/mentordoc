import {IRootState} from "./model/root-state";

type Selector = (state: IRootState) => {
    [key: string]: any,
};

interface IPropMap {
    selector: {
        [key: string]: any;
    };
}

export function combineSelectors<T>(...selectors: Selector[]): (state: IRootState) => IPropMap {
    return (state: IRootState): IPropMap => {
        const map: IPropMap = {
            selector: {

            },
        };

        selectors.map((selector: Selector): { [key: string]: any } => {
            return selector(state);
        });

        return map;
    };
}
