import {IRootState} from "./model/root-state";
import {Dispatch} from "react";
import {AnyAction} from "redux";

type Selector = (state: IRootState) => {
    [key: string]: any
};

type PropMap = {
    selector: {
        [key: string]: any;
    }
};

export function combineSelectors<T>(...selectors: Selector[]): (state: IRootState) => PropMap {
    return (state: IRootState): PropMap => {
        const map: PropMap = {
            selector: {

            }
        };

        selectors.map((selector: Selector): { [key: string]: any } =>{
            return selector(state)
        });

        return map;
    }
}