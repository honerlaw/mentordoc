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

interface ISelectorPropMap<T = IDefaultPropMap> {
    selector: T;
}

interface IDispatchPropMap<T = IDefaultPropMap> {
    dispatch: T;
}

export function ConnectProps(
    mapStateToProps?: MapToProps,
    mapDispatchToProps?: MapToProps,
): ClassDecorator {
    return connect(mapStateToProps, mapDispatchToProps) as any;
}

export function CombineSelectors(...selectors: Selector[]): (state: IRootState) => ISelectorPropMap {
    return (state: IRootState): ISelectorPropMap => {
        const map: ISelectorPropMap = {
            selector: {},
        };

        selectors.map((selector: Selector): { [key: string]: any } => {
            return selector(state);
        });

        return map;
    };
}

export function CombineDispatchers(...dispatchers: Dispatcher[]): (dispatch: Dispatch<AnyAction>) => IDispatchPropMap {
    return (dispatch: Dispatch<AnyAction>): IDispatchPropMap => {
        const map: IDispatchPropMap = {
            dispatch: {},
        };

        dispatchers.map((dispatcher: Dispatcher): { [key: string]: any } => {
            return dispatcher(dispatch);
        });

        return map;
    };
}
