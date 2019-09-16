import * as React from "react";
import "./dashboard-viewer.scss";
import {DocumentRenderer} from "./dashboard-viewer/document-renderer";
import {WithRouter} from "@honerlawd/mentordoc-frontend-shared/dist/store/decorator/with-router";
import {RouteComponentProps} from "react-router";

interface IRouteProps {
    orgId: string;
    docId: string;
}

interface IProps extends Partial<RouteComponentProps<IRouteProps>> {

}

@WithRouter()
export class DashboardViewer extends React.PureComponent<IProps, {}> {

    public render(): JSX.Element {
        return <div className={"dashboard-viewer"}>
            {this.renderViewer()}
        </div>
    }

    private renderViewer(): JSX.Element | null {
        if (this.props.match!.params.docId && this.props.match!.params.orgId) {
            return <DocumentRenderer/>;
        }
        return null;
    }

}
