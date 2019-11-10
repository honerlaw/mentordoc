import * as React from "react";
import {Page} from "../shared/page";
import "./organization-page.scss";
import {onChangeSetState} from "../../util";
import {InviteModal} from "./invite/invite-modal";
import {
    CombineSelectors,
    ConnectProps,
    ISelectorPropMap
} from "@honerlawd/mentordoc-frontend-shared/dist/store/decorator/connect-props";
import {
    ISetCurrentOrganizationSelector,
    SetCurrentOrganization
} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/organization/set-current-organization";

type Props = Partial<ISelectorPropMap<ISetCurrentOrganizationSelector>>;

interface IState {
    title: string;
    isInviteModalVisible: boolean;
}

@ConnectProps(CombineSelectors(SetCurrentOrganization.selector))
export class OrganizationPage extends React.PureComponent<Props, IState> {

    public constructor(props: Props) {
        super(props);

        this.state = {
            title: props.selector!.currentOrganization!.model.name,
            isInviteModalVisible: false
        };

        this.onInviteModalRequestClose = this.onInviteModalRequestClose.bind(this);
        this.onInviteClick = this.onInviteClick.bind(this);
    }

    public render(): JSX.Element {
        return <Page className={"organization-page"}>
            <div className={"settings"}>

                <h3>Organization Settings</h3>

                <input type={"text"}
                       placeholder={"organization name"}
                       value={this.state.title}
                       onChange={onChangeSetState<IState>("title", this)}/>

                <section className={"user-management"}>
                    <header className={"user-management-header"}>
                        <h5>Organization Users</h5>
                        <button onClick={this.onInviteClick}>Invite User</button>
                    </header>
                    <table>
                        <thead>
                        <tr>
                            <th>Person</th>
                            <th>Permission</th>
                            <th>Status</th>
                            <th>Options</th>
                        </tr>
                        </thead>
                    </table>
                </section>

                <InviteModal isVisible={this.state.isInviteModalVisible}
                             onRequestClose={this.onInviteModalRequestClose}/>

            </div>
        </Page>;
    }

    private onInviteModalRequestClose(): void {
        this.setState({
            isInviteModalVisible: false
        });
    }

    private onInviteClick(): void {
        this.setState({
            isInviteModalVisible: true
        });
    }

}