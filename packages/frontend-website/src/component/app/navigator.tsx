import * as React from "react";
import {
    CombineDispatchers, CombineSelectors, ConnectProps,
    IDispatchPropMap,
    ISelectorPropMap
} from "@honerlawd/mentordoc-frontend-shared/dist/store/decorator/connect-props";
import {
    ISetOrganizationsSelector,
    SetOrganizations
} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/organization/set-organizations";
import {
    FetchOrganizations,
    IFetchOrganizationsDispatch
} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/organization/fetch-organizations";
import "./navigator.scss";
import {AclOrganization} from "@honerlawd/mentordoc-frontend-shared/dist/store/model/organization/acl-organization";
import {NavigatorItem} from "./navigator/navigator-item";

interface IProps extends Partial<ISelectorPropMap<ISetOrganizationsSelector> & IDispatchPropMap<IFetchOrganizationsDispatch>> {

}

@ConnectProps(CombineSelectors(SetOrganizations.selector), CombineDispatchers(FetchOrganizations.dispatch))
export class Navigator extends React.PureComponent<IProps, {}> {

    public async componentWillMount(): Promise<void> {
        await this.props.dispatch!.fetchOrganizations();
    }

    public render(): JSX.Element {
        return <div className={"dashboard-navigator"}>
            <h4>Navigation</h4>
            {this.props.selector!.organizations.map((org: AclOrganization): JSX.Element => {
                return <NavigatorItem key={org.model.id} item={org} />;
            })}
        </div>;
    }

}