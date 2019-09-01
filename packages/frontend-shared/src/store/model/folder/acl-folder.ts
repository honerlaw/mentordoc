import {Folder} from "./folder";
import {Type, Exclude, Expose} from "class-transformer";

@Exclude()
export class AclFolder {

    @Expose()
    @Type(() => Folder)
    public model: Folder;

    @Expose()
    public actions: string[];

    public hasAction(action: string): boolean {
        return this.actions.indexOf(action) !== -1;
    }

}