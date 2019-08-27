import {Exclude, Expose} from "class-transformer";

@Exclude()
export class AuthenticationData {

    @Expose()
    public accessToken: string;

    @Expose()
    public refreshToken: string;

}
