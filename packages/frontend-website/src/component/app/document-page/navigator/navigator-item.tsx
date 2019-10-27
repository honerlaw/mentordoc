import * as React from "react";
import {
    AclOrganization,
    isAclOrganization
} from "@honerlawd/mentordoc-frontend-shared/dist/store/model/organization/acl-organization";
import {OrganizationItemView} from "./organization-item-view";
import {AclFolder, isAclFolder} from "@honerlawd/mentordoc-frontend-shared/dist/store/model/folder/acl-folder";
import {FolderItemView} from "./folder-item-view";
import {AclDocument, isAclDocument} from "@honerlawd/mentordoc-frontend-shared/dist/store/model/document/acl-document";
import {DocumentItemView} from "./document-item-view";

interface IProps {
    item: AclOrganization | AclFolder | AclDocument;
}

export class NavigatorItem extends React.PureComponent<IProps, {}> {

    public render(): JSX.Element | null {
        if (isAclOrganization(this.props.item)) {
            return <OrganizationItemView organization={this.props.item} />;
        }
        if (isAclFolder(this.props.item)) {
            return <FolderItemView folder={this.props.item} />;
        }
        if (isAclDocument(this.props.item)) {
            return <DocumentItemView document={this.props.item}/>;
        }
        return null;
    }

}