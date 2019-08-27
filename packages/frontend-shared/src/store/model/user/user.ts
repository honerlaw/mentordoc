import {Exclude, Expose} from "class-transformer";

@Exclude()
export class User extends Entity {

    @Expose()
    public email: string;

}
