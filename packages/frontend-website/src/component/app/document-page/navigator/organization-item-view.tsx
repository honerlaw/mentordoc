import * as React from "react";
import {AclOrganization} from "@honerlawd/mentordoc-frontend-shared/dist/store/model/organization/acl-organization";
import {DropdownButton, IDropdownButtonOption} from "../../../shared/dropdown-button";
import {
    CombineDispatchers, CombineSelectors,
    ConnectProps, IDispatchPropMap, ISelectorPropMap
} from "@honerlawd/mentordoc-frontend-shared/dist/store/decorator/connect-props";
import {
    FetchFolders,
    IFetchFoldersDispatch
} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/folder/fetch-folders";
import {
    ISetFoldersSelector, SelectorSetFoldersChildValue,
    SetFolders
} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/folder/set-folders";
import {AclFolder} from "@honerlawd/mentordoc-frontend-shared/dist/store/model/folder/acl-folder";
import {NavigatorItem} from "./navigator-item";
import {
    CreateFolder,
    ICreateFolderDispatch
} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/folder/create-folder";
import {
    CreateDocument,
    ICreateDocumentDispatch
} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/document/create-document";
import {
    ISetDocumentsSelector, SelectorSetDocumentsChildValue,
    SetDocuments
} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/document/set-documents";
import {
    FetchDocuments,
    IFetchDocumentsDispatch
} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/document/fetch-documents";
import {AclDocument} from "@honerlawd/mentordoc-frontend-shared/dist/store/model/document/acl-document";
import * as icon from "../../../../../images/ellipsis.svg";
import "./organization-item-view.scss";

interface IProps extends Partial<IDispatchPropMap<IFetchFoldersDispatch & ICreateFolderDispatch & IFetchDocumentsDispatch & ICreateDocumentDispatch>
    & ISelectorPropMap<ISetFoldersSelector & ISetDocumentsSelector>> {
    organization: AclOrganization;
}

@ConnectProps(
    CombineSelectors(SetFolders.selector, SetDocuments.selector),
    CombineDispatchers(FetchFolders.dispatch, CreateFolder.dispatch, FetchDocuments.dispatch, CreateDocument.dispatch)
)
export class OrganizationItemView extends React.PureComponent<IProps, {}> {

    public constructor(props: IProps) {
        super(props);

        this.createFolder = this.createFolder.bind(this);
        this.createDocument = this.createDocument.bind(this);
    }

    public async componentDidMount(): Promise<void> {
        await this.props.dispatch!.fetchFolders({
            organizationId: this.props.organization.model.id,
            parentFolderId: null
        });

        await this.props.dispatch!.fetchDocuments({
            organizationId: this.props.organization.model.id,
            folderId: null
        });
    }

    public async componentDidUpdate(prevProps: Readonly<IProps>, prevState: Readonly<{}>, snapshot?: any): Promise<void> {
        if (prevProps.organization.model.id !== this.props.organization.model.id) {
            await this.componentDidMount();
        }
    }

    public render(): JSX.Element {
        return <div className={"organization-item-view"}>

            <div className={"organization-item-view-header"}>
                <h5>{this.props.organization.model.name}</h5>
                <DropdownButton icon={icon} position={"left"} options={this.getOptions()}/>
            </div>

            {this.renderFolders()}
            {this.renderDocuments()}
        </div>;
    }

    private renderFolders(): JSX.Element[] | null {
        const selector: SelectorSetFoldersChildValue = this.props.selector!.getFolders("child") as SelectorSetFoldersChildValue;
        const folders: AclFolder[] = selector(this.props.organization.model.id, null);

        if (folders.length === 0) {
            return null;
        }

        return folders.map((folder: AclFolder): JSX.Element => {
            return <NavigatorItem key={folder.model.id} item={folder} />;
        });
    }

    private renderDocuments(): JSX.Element[] | null {
        const selector: SelectorSetDocumentsChildValue = this.props.selector!.getDocuments("child") as SelectorSetDocumentsChildValue;
        const documents: AclDocument[] = selector(this.props.organization.model.id, null);

        if (documents.length === 0) {
            return null;
        }

        return documents.map((doc: AclDocument): JSX.Element => {
            return <NavigatorItem key={doc.model.id} item={doc} />;
        });
    }

    private getOptions(): IDropdownButtonOption[] {
        const options: IDropdownButtonOption[] = [];

        if (this.props.organization.hasAction("create:folder")) {
            options.push({
                label: "add folder",
                onClick: this.createFolder
            });
        }

        if (this.props.organization.hasAction("create:document")) {
            options.push({
                label: "add document",
                onClick: this.createDocument
            });
        }

        return options;
    }

    private async createFolder(): Promise<void> {
        await this.props.dispatch!.createFolder({
            organizationId: this.props.organization.model.id,
            parentFolderId: null,
            name: "New Folder"
        });
    }

    private async createDocument(): Promise<void> {
        await this.props.dispatch!.createDocument({
            organizationId: this.props.organization.model.id,
            folderId: null,
            name: "New Document",
            content: "testing"
        });
    }

}
