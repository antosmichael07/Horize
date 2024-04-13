document.getElementById('login_button').addEventListener('click', function () {
    document.getElementById('login_button').innerHTML = "<span class=\"spinner-border\" role=\"status\" aria-hidden=\"true\"></span>"
    document.getElementById('login_button').disabled = true

    fetch("/post_login", {
        method: "POST",
        body: document.getElementById('login_login').value + "|" +
            document.getElementById('login_password').value
    }).then((response) => {
        response.text().then(function (data) {
            document.getElementById('login_button').disabled = false
            document.getElementById('login_button').innerHTML = "login"

            var res = data.split('\\')

            document.getElementById('login_button').classList.remove('my-3')
            document.getElementById('login_response').classList.remove('text-success')
            document.getElementById('login_response').classList.remove('text-danger')
            document.getElementById('login_response').classList.add(res[0])

            document.getElementById('login_response').innerHTML = res[1]
            document.getElementById('login_response').style.display = "block"

            if (res[0] == "text-success") {
                localStorage.setItem("account_token", res[2])
                window.location.href = "/play"
            }
        });
    }).catch((error) => {
        console.log(error)
    });
})
