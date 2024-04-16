function logout() {
    fetch("/post_logout", {
        method: "POST",
        body: localStorage.getItem("account_token")
    }).then((response) => {
        response.text().then(function (data) {
            var res = data.split('\\')

            if (res[0] == "text-success") {
                localStorage.removeItem("account_token")
                window.location.href = "/login"
            } else {
                alert("Error: " + res[1])
            }
        });
    }).catch((error) => {
        console.log(error)
    });
}
