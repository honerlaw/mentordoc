import * as React from "react";
import {
    CombineSelectors,
    ConnectProps,
    ISelectorPropMap
} from "@honerlawd/mentordoc-frontend-shared/dist/store/decorator/connect-props";
import {AddAlert, IAddAlertSelector} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/alert/add-alert";
import {Alert} from "@honerlawd/mentordoc-frontend-shared/dist/store/model/alert/alert";
import "./alert-list.scss";
import {AlertListItem} from "./alert-list-item";

interface IProps extends Partial<ISelectorPropMap<IAddAlertSelector>> {
    target?: string;
}

@ConnectProps(CombineSelectors(AddAlert.selector))
export class AlertList extends React.PureComponent<IProps, {}> {

    public render(): JSX.Element {
        return <div className={"alert-list"}>
            {this.renderAlerts()}
        </div>;
    }

    private renderAlerts(): JSX.Element[] {
        if (!this.props.selector!.alerts) {
            return [];
        }

        const alerts: JSX.Element[] = [];
        this.props.selector!.alerts.forEach((alert: Alert): void => {
            if (alert.target !== this.props.target) {
                return;
            }

            alerts.push(<AlertListItem key={alert.getKey()} alert={alert} />);
        });

        return alerts;
    }

}