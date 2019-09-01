import {connect, MapDispatchToPropsNonObject, MapStateToPropsParam} from "react-redux";
import {AnyAction, Dispatch} from "redux";
import {IRootState} from "../model/root-state";

type MapToProps = MapStateToPropsParam<any, any, any> | MapDispatchToPropsNonObject<any, any> | undefined;

type Selector = (state: IRootState) => {
    [key: string]: any,
};

type Dispatcher = (dispatch: Dispatch<any>) => {
    [key: string]: any,
};

interface IDefaultPropMap {
    [key: string]: any;
}

export interface ISelectorPropMap<T = IDefaultPropMap> {
    selector: T;
}

export interface IDispatchPropMap<T = IDefaultPropMap> {
    dispatch: T;
}

interface IMap {
    [key: string]: any;
}

export function ConnectProps(
    mapStateToProps?: MapToProps,
    mapDispatchToProps?: MapToProps,
): ClassDecorator {
    return connect(mapStateToProps, mapDispatchToProps, null, {pure: false}) as any;
}

export function CombineSelectors(...selectors: Selector[]): (state: IRootState) => ISelectorPropMap {
    return (state: IRootState): ISelectorPropMap => {
        return {
            selector: selectors.map((selector: Selector): IMap => {
                return selector(state);
            }).reduce((one: IMap, two: IMap): IMap => {
                return {
                    ...one,
                    ...two,
                }
            }, {})
        };
    };
}

export function CombineDispatchers(...dispatchers: Dispatcher[]): (dispatch: Dispatch<AnyAction>) => IDispatchPropMap {
    return (dispatch: Dispatch<AnyAction>): IDispatchPropMap => {
        return {
            dispatch: dispatchers.map((dispatcher: Dispatcher): IMap => {
                return dispatcher(dispatch);
            }).reduce((one: IMap, two: IMap): IMap => {
                return {
                    ...one,
                    ...two,
                }
            }, {})
        };
    };
}
