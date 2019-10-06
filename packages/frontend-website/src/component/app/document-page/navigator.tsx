import * as React from "react";
import {
    CombineSelectors, ConnectProps,
    ISelectorPropMap
} from "@honerlawd/mentordoc-frontend-shared/dist/store/decorator/connect-props";
import "./navigator.scss";
import {AclOrganization} from "@honerlawd/mentordoc-frontend-shared/dist/store/model/organization/acl-organization";
import {NavigatorItem} from "./navigator/navigator-item";
import {
    ISetCurrentOrganizationSelector,
    SetCurrentOrganization
} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/organization/set-current-organization";

interface IProps extends Partial<ISelectorPropMap<ISetCurrentOrganizationSelector>> {

}

@ConnectProps(CombineSelectors(SetCurrentOrganization.selector))
export class Navigator extends React.PureComponent<IProps, {}> {

    public render(): JSX.Element {
        return <div className={"dashboard-navigator"}>
            <div className={"dashboard-empty-spacer"} />
            <div className={"dashboard-navigator-container"}>
                {this.renderOrganization()}
            </div>
        </div>;
    }

    private renderOrganization(): JSX.Element | null {
        const org: AclOrganization | null = this.props.selector!.currentOrganization;
        if (!org) {
            return null;
        }
        return <NavigatorItem item={org} key={org.model.id} />;
    }

}