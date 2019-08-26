export function onChangeSetState<State>(key: keyof State, obj: any): (event: React.ChangeEvent<HTMLInputElement>) => void {
    return (event: React.ChangeEvent<HTMLInputElement>): void => {
        obj.setState({
            [key]: event.target.value
        });
    }
}
