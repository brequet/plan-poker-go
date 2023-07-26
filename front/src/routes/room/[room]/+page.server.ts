export const load = ({cookies}) => {
    const nickname = cookies.get('nickname');

    // todo: fetch room info (code, name)

    return {
        nickname
    };
}