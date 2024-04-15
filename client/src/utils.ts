export function isLoggedIn(): boolean {
    return document.cookie.split(';').find((cookie) => cookie === "loggedIn=1") !== undefined;
}