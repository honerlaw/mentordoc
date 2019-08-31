import {Exclude, Expose} from "class-transformer";
import {Entity} from "../entity";

@Exclude()
export class User extends Entity {

    @Expose()
    public email: string;

}
