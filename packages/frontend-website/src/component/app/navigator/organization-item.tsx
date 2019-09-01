import * as React from "react";
import {AclOrganization} from "@honerlawd/mentordoc-frontend-shared/dist/store/model/organization/acl-organization";
import * as icon from "../../../../images/organization.svg";

interface IProps {
    organization: AclOrganization;
}

export class OrganizationItem extends React.PureComponent<IProps, {}> {

    public render(): JSX.Element {
        return <div className={"navigator-item organization-item"}>
            <img src={icon} alt={"organization"} />
            <span>
                {this.props.organization.model.name}
            </span>
        </div>;
    }

}