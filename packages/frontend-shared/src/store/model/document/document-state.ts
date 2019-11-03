import {AclDocument} from "./acl-document";
import {AclOrganization} from "../organization/acl-organization";
import {AclFolder} from "../folder/acl-folder";
import {ISetSearchDocuments} from "../../action/document/set-search-documents";

export type DocumentPath = Array<AclOrganization | AclFolder | AclDocument>;

export interface IDocumentState extends ISetSearchDocuments {
    documentMap: Record<string, AclDocument[]>;
    fullDocument: AclDocument | null;
    documentPath: DocumentPath;
}

export const INITIAL_DOCUMENT_STATE: IDocumentState = {
    documentMap: {},
    fullDocument: null,
    documentPath: [],
    searchDocuments: null
};
