import * as React from "react";
import {
    CombineSelectors,
    ConnectProps,
    ISelectorPropMap
} from "@honerlawd/mentordoc-frontend-shared/dist/store/decorator/connect-props";
import {AddAlert, IAddAlertSelector} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/alert/add-alert";
import {Alert} from "@honerlawd/mentordoc-frontend-shared/dist/store/model/alert/alert";

interface IProps extends Partial<ISelectorPropMap<IAddAlertSelector>> {

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

        return this.props.selector!.alerts.map((alert: Alert): JSX.Element => {
            return <div key={alert.getKey()}>
                {alert.message}
            </div>;
        });
    }

}