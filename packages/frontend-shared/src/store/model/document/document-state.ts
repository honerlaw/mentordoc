import {AclDocument} from "./acl-document";
import {Organization} from "../organization/organization";
import {Folder} from "../folder/folder";

export type DocumentPath = Array<Organization | Folder | Document>;

export interface IDocumentState {
    documentMap: Record<string, AclDocument[]>;
    fullDocument: AclDocument | null;
    documentPath: DocumentPath;
}

export const INITIAL_DOCUMENT_STATE: IDocumentState = {
    documentMap: {},
    fullDocument: null,
    documentPath: []
};
