import * as React from "react";
import {Alert} from "@honerlawd/mentordoc-frontend-shared/dist/store/model/alert/alert";
import {
    CombineDispatchers,
    ConnectProps, IDispatchPropMap
} from "@honerlawd/mentordoc-frontend-shared/dist/store/decorator/connect-props";
import {
    IRemoveAlertDispatch,
    RemoveAlert
} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/alert/remove-alert";

interface IProps extends Partial<IDispatchPropMap<IRemoveAlertDispatch>> {
    alert: Alert;
}

@ConnectProps(null, CombineDispatchers(RemoveAlert.dispatch))
export class AlertListItem extends React.PureComponent<IProps, {}> {

    private timeout: any;

    public constructor(props: IProps) {
        super(props);

        this.remove = this.remove.bind(this);
    }

    public componentDidMount(): void {
        if (this.props.alert.lifespan) {
            this.timeout = setTimeout(this.remove, this.props.alert.lifespan);
        }
    }

    public componentWillUnmount(): void {
        this.remove();
    }

    public render(): JSX.Element {
        return <div className={`alert-list-item ${this.props.alert.type}`} key={this.props.alert.getKey()}>
            <span>{this.props.alert.message}</span>
            <button onClick={this.remove}>&#10005;</button>
        </div>;
    }

    private remove(): void {
        if (this.timeout) {
            clearTimeout(this.timeout);
        }

        this.props.dispatch!.removeAlert({
            alert: this.props.alert
        });
    }

}