import * as React from "react";
import {NavigatorItemView} from "./navigator-item-view";
import {IDropdownButtonOption} from "../../../shared/dropdown-button";
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
    ISetDocumentsSelector, SelectorSetDocumentsChildValue,
    SetDocuments
} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/document/set-documents";
import {
    FetchDocuments,
    IFetchDocumentsDispatch
} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/document/fetch-documents";
import {
    CreateDocument,
    ICreateDocumentDispatch
} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/document/create-document";
import {AclDocument} from "@honerlawd/mentordoc-frontend-shared/dist/store/model/document/acl-document";
import {
    DeleteFolder,
    IDeleteFolderDispatch
} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/folder/delete-folder";

interface IProps extends Partial<IDispatchPropMap<IFetchFoldersDispatch & ICreateFolderDispatch & IFetchDocumentsDispatch & ICreateDocumentDispatch & IDeleteFolderDispatch>
    & ISelectorPropMap<ISetFoldersSelector & ISetDocumentsSelector>> {
    folder: AclFolder;
}

@ConnectProps(
    CombineSelectors(SetFolders.selector, SetDocuments.selector),
    CombineDispatchers(FetchFolders.dispatch, CreateFolder.dispatch, FetchDocuments.dispatch, CreateDocument.dispatch, DeleteFolder.dispatch)
)
export class FolderItemView extends React.PureComponent<IProps, {}> {

    public constructor(props: IProps) {
        super(props);

        this.state = {
            isExpanded: false
        };

        this.createFolder = this.createFolder.bind(this);
        this.createDocument = this.createDocument.bind(this);
        this.deleteFolder = this.deleteFolder.bind(this);
        this.onExpand = this.onExpand.bind(this);
    }

    public render(): JSX.Element {
        return <NavigatorItemView title={this.props.folder.model.name}
                                  hasChildren={this.props.folder.model.childCount > 0}
                                  isExpanded={false}
                                  onExpand={this.onExpand}
                                  options={this.getOptions()}>
            {this.renderFolders()}
            {this.renderDocuments()}
        </NavigatorItemView>
    }

    private renderFolders(): JSX.Element[] | null {
        const selector: SelectorSetFoldersChildValue = this.props.selector!.getFolders("child") as SelectorSetFoldersChildValue;
        const folders: AclFolder[] = selector(this.props.folder.model.organizationId, this.props.folder.model.id);

        if (folders.length === 0) {
            return null;
        }

        return folders.map((folder: AclFolder): JSX.Element => {
            return <NavigatorItem key={folder.model.id} item={folder}/>;
        });
    }

    private renderDocuments(): JSX.Element[] | null {
        const selector: SelectorSetDocumentsChildValue = this.props.selector!.getDocuments("child") as SelectorSetDocumentsChildValue;
        const documents: AclDocument[] = selector(this.props.folder.model.organizationId, this.props.folder.model.id);

        if (documents.length === 0) {
            return null;
        }

        return documents.map((doc: AclDocument): JSX.Element => {
            return <NavigatorItem key={doc.model.id} item={doc} />;
        });
    }

    private getOptions(): IDropdownButtonOption[] {
        const options: IDropdownButtonOption[] = [];

        if (this.props.folder.hasAction("create:folder")) {
            options.push({
                label: "add folder",
                onClick: this.createFolder
            });
        }

        if (this.props.folder.hasAction("create:document")) {
            options.push({
                label: "add document",
                onClick: this.createDocument
            });
        }

        if (this.props.folder.hasAction("delete") && this.props.folder.model.childCount == 0) {
            options.push({
                label: "delete",
                onClick: this.deleteFolder
            });
        }

        return options;
    }

    private async createFolder(): Promise<void> {
        await this.props.dispatch!.createFolder({
            organizationId: this.props.folder.model.organizationId,
            parentFolderId: this.props.folder.model.id,
            name: "New Folder"
        });
    }

    private async createDocument(): Promise<void> {
        await this.props.dispatch!.createDocument({
            organizationId: this.props.folder.model.organizationId,
            folderId: this.props.folder.model.id,
            name: "New Document",
            content: "testing"
        });
    }

    private async deleteFolder(): Promise<void> {
        await this.props.dispatch!.deleteFolder({
            folderId: this.props.folder.model.id
        });
    }

    private async onExpand(): Promise<void> {
        await this.props.dispatch!.fetchFolders({
            organizationId: this.props.folder.model.organizationId,
            parentFolderId: this.props.folder.model.id
        });

        await this.props.dispatch!.fetchDocuments({
            organizationId: this.props.folder.model.organizationId,
            folderId: this.props.folder.model.id
        });
    }

}
