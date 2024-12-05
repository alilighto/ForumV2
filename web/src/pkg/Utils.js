const showError = (status, message) => {
    const app = document.querySelector('#app')

    const titles = {
        400: "400 Bad Request",
        401: "401 Unauthorized",
        403: "403 Forbidden",
        404: "404 Not Found",
        405: "405 Method Not Allowed",
        500: "500 Internal Server Error",
        503: "503 Service Unavailable"
    }

    app.innerHTML = `
        <div class="errorDiv">
        <h1>${titles[status]}</h1><br>
        <h2>${message || ''}</h2>
        </div>
    `
}

const getUser = () => {
    return {
        id: localStorage.getItem('id'),
        role: localStorage.getItem('role'),
        token: localStorage.getItem('token')
    }
}

const logOut = () => {
    localStorage.removeItem('id')
    localStorage.removeItem('role')
    localStorage.removeItem('token')
}

const parseJwt = (token) => {
    var base64Url = token.split('.')[1];
    var base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    var jsonPayload = decodeURIComponent(atob(base64).split('').map(function (c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''));

    return JSON.parse(jsonPayload);
};

export default {showError, getUser, logOut, parseJwt};