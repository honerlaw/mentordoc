import * as React from "react";
import {Modal} from "../../shared/modal";
import "./invite-modal.scss";
import {DropdownButton} from "../../shared/dropdown-button";

interface IProps {
    isVisible: boolean;
    onRequestClose: () => void;
}

export class InviteModal extends React.PureComponent<IProps, {}> {

    public constructor(props: IProps) {
        super(props);

        this.onSubmit = this.onSubmit.bind(this);
    }


    public render(): JSX.Element {
        return <Modal isVisible={this.props.isVisible} onRequestClose={this.props.onRequestClose}>

            <form className={"invite-modal"} onSubmit={this.onSubmit}>
                <h3>Invite Users</h3>
                <input type={"text"} placeholder={"emails"}/>
                <div className={"options"}>

                    <DropdownButton label={"select permission"} showSelected={true} options={[{
                        label: "contributor",
                        onClick: () => {}
                    }, {
                        label: "read only",
                        onClick: () => {}
                    }]} />
                    <button className={"send-button"}>send</button>
                </div>
            </form>

        </Modal>;
    }

    private onSubmit(): void {
        // @todo this should actually send off / email the users
    }

}
