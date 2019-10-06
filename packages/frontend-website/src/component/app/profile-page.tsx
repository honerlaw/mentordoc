import * as React from "react";
import {Page} from "../shared/page";
import "./profile-page.scss";

export class ProfilePage extends React.PureComponent<{}, {}> {

    public render(): JSX.Element {
        return <Page>
            <div id={"profile-page"}>
                <div className={"settings"}>
                    profile
                </div>
            </div>
        </Page>;
    }

}