import * as React from "react";
import {WithRouter} from "@honerlawd/mentordoc-frontend-shared/dist/store/decorator/with-router";
import {
    CombineDispatchers,
    CombineSelectors,
    ConnectProps, IDispatchPropMap, ISelectorPropMap
} from "@honerlawd/mentordoc-frontend-shared/dist/store/decorator/connect-props";
import {
    FetchFullDocument,
    IFetchFullDocumentDispatch
} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/document/fetch-full-document";
import {
    ISetFullDocumentDispatch, ISetFullDocumentSelector,
    SetFullDocument
} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/document/set-full-document";
import {RouteComponentProps} from "react-router";
import {AclDocument} from "@honerlawd/mentordoc-frontend-shared/dist/store/model/document/acl-document";
import {DocumentViewer} from "./document-viewer";
import {DocumentEditor} from "./document-editor";
import {
    ISetDocumentsSelector, SelectorSetDocumentsMapValue,
    SetDocuments
} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/document/set-documents";
import {
    ISetFoldersSelector,
    SelectorSetFoldersMapValue, SetFolders
} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/folder/set-folders";
import {AclFolder} from "@honerlawd/mentordoc-frontend-shared/dist/store/model/folder/acl-folder";
import {
    ISetOrganizationsSelector,
    SetOrganizations
} from "@honerlawd/mentordoc-frontend-shared/dist/store/action/organization/set-organizations";
import {AclOrganization} from "@honerlawd/mentordoc-frontend-shared/dist/store/model/organization/acl-organization";

export interface IRouteProps {
    orgId: string;
    docId: string;
}

interface IProps extends Partial<IDispatchPropMap<IFetchFullDocumentDispatch & ISetFullDocumentDispatch> &
    ISelectorPropMap<ISetFullDocumentSelector & ISetDocumentsSelector & ISetFoldersSelector & ISetOrganizationsSelector> &
    RouteComponentProps<IRouteProps>> {
}

interface IState {
    isEditing: boolean;
}

@WithRouter()
@ConnectProps(
    CombineSelectors(SetFullDocument.selector, SetDocuments.selector, SetFolders.selector, SetOrganizations.selector),
    CombineDispatchers(FetchFullDocument.dispatch, SetFullDocument.dispatch)
)
export class DocumentRenderer extends React.PureComponent<IProps, IState> {

    public constructor(props: IProps) {
        super(props);

        this.state = {
            isEditing: false
        };
    }

    public async componentWillReceiveProps(nextProps: Readonly<IProps>, nextContext: any): Promise<void> {
        const didOrgChange: boolean = this.props.match!.params.orgId !== nextProps.match!.params.orgId;
        const didDocChange: boolean = this.props.match!.params.docId !== nextProps.match!.params.docId;
        if(didOrgChange || didDocChange) {
            await this.componentWillMount(nextProps.match!.params.docId);
        }
    }

    public async componentWillMount(docId: string = this.props.match!.params.docId): Promise<void> {
        this.props.dispatch!.setFullDocument({
            fullDocument: null
        });

        await this.props.dispatch!.fetchFullDocument({
            documentId: docId
        });
    }

    public render(): JSX.Element | null {
        const doc: AclDocument | null = this.props.selector!.fullDocument;
        if (!doc) {
            return null;
        }

        if (this.state.isEditing) {
            return <DocumentEditor document={doc}/>;
        }
        return <div>
            <div className={"document-action-bar"}>
                <span>{this.getPath(doc).join(" > ")}</span>
                <div className={"actions"}>
                    <button>modifiy</button>
                </div>
            </div>
            <DocumentViewer document={doc}/>
        </div>;
    }

    // @todo this should be moved to the backend
    private getPath(doc: AclDocument): string[] {
        if (doc.model.folderId === null) {
            // @todo find the organization name
            return [this.getOrganizationName(doc.model.organizationId), doc.model.drafts[0].name]
        }

        const folderMap: SelectorSetFoldersMapValue = this.props.selector!.getFolders("map") as SelectorSetFoldersMapValue;
        const folderArr: AclFolder[][] = Object.values(folderMap);

        const path: string[] = this.recursivelyFindFolderPath(doc.model.folderId, folderArr);

        path.reverse();

        path.push(doc.model.drafts[0].name);

        return path;
    }

    private getOrganizationName(orgId: string): string {
        const orgs: AclOrganization[] | null = this.props.selector!.organizations;
        if (!orgs) {
            return "Unknown";
        }

        const org: AclOrganization | undefined = orgs.find((org: AclOrganization): boolean => org.model.id === orgId);
        if (!org) {
            return "Unknown";
        }

        return org.model.name;
    }

    private recursivelyFindFolderPath(folderId: string, folderArr: AclFolder[][]): string[] {
        let path: string[] = [];
        for (let i: number = 0; i < folderArr.length; ++i) {
            const folders: AclFolder[] = folderArr[i];
            for (let j: number = 0; j < folders.length; ++j) {
                const folder: AclFolder = folders[j];

                // found the current folder
                if (folder.model.id === folderId) {
                    path.push(folder.model.name);
                    if (folder.model.parentFolderId !== null) {
                        path = path.concat(this.recursivelyFindFolderPath(folder.model.parentFolderId, folderArr))
                    } else {
                        path.push(this.getOrganizationName(folder.model.organizationId));
                    }
                }
            }
        }
        return path;
    }

}
