import {AclDocument} from "./acl-document";
import {AclOrganization} from "../organization/acl-organization";
import {AclFolder} from "../folder/acl-folder";

export type DocumentPath = Array<AclOrganization | AclFolder | AclDocument>;

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
