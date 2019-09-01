import {Organization} from "./organization";
import {Exclude, Expose, Type} from "class-transformer";

@Exclude()
export class AclOrganization {

    @Expose()
    @Type(() => Organization)
    public model: Organization;

    @Expose()
    public actions: string[];

    public hasAction(action: string): boolean {
        return this.actions.indexOf(action) !== -1;
    }

}