import * as React from "react";
import "./dashboard-viewer.scss";
import {DocumentViewer} from "./dashboard-viewer/document-viewer";

export class DashboardViewer extends React.PureComponent<{}, {}> {

    public render(): JSX.Element {
        return <div className={"dashboard-viewer"}>
            <DocumentViewer />
        </div>
    }

}
