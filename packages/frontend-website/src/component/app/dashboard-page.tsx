import * as React from "react";
import {Page} from "../shared/page";
import {Navigator} from "./dashboard-page/navigator";
import {DashboardViewer} from "./dashboard-page/dashboard-viewer";
import "./dashboard-page.scss";
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
export class DashboardPage extends React.PureComponent<IProps, {}> {

    public async componentWillMount(): Promise<void> {
        if (!this.props.match.params.orgId) {
            if (this.props.selector.organizations) {
                this.props.history.push(`/app/${this.props.selector.organizations[0].model.id}`)
            }
        }
    }

    public render(): JSX.Element {
        return <Page className={"dashboard-page"}>
            <Navigator/>
            <DashboardViewer/>
        </Page>;
    }

}
