import * as React from "react";
import {Page} from "../shared/page";
import {Navigator} from "./navigator";
import {DashboardViewer} from "./dashboard-viewer";
import "./dashboard.scss";
import {RouteComponentProps} from "react-router";
import {
    CombineSelectors,
    ConnectProps, ISelectorPropMap
} from "@honerlawd/mentordoc-frontend-shared/dist/store/decorator/connect-props";
import {
    ISetOrganizationsSelector,
    SetOrganizations
} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/organization/set-organizations";

interface IProps extends RouteComponentProps<IRouteParams>, ISelectorPropMap<ISetOrganizationsSelector> {

}

interface IRouteParams {
    orgId: string;
    docId: string;
}

@ConnectProps(CombineSelectors(SetOrganizations.selector))
export class Dashboard extends React.PureComponent<IProps, {}> {

    public async componentWillMount(): Promise<void> {
        if (!this.props.match.params.orgId) {
            if (this.props.selector.organizations) {
                this.props.history.push(`/app/${this.props.selector.organizations[0].model.id}`)
            }
        }
    }

    public render(): JSX.Element {
        return <Page>
            <div id={"dashboard-container"}>
                <Navigator/>
                <DashboardViewer/>
            </div>
        </Page>;
    }

}
