function checkLoginValidity(str) {
    var pipeCount = 0
    for (var i = 0; i < str.length; i++) {
        if (str[i] == '|') {
            pipeCount++
        }
    }
    if (pipeCount != 1) {
        return "text-danger\\Invalid data (cannot contain | or the data sent was too long)"
    }

    var data = str.split('|')

    if (data[0].replaceAll(" ", "") == "") {
        return "text-danger\\Invalid login (cannot be empty)"
    }

    if (!/^[a-zA-Z0-9]*$/.test(data[0])) {
        return "text-danger\\Invalid login (can only contain letters and numbers)"
    }

    if (data[0].length > 16) {
        return "text-danger\\Invalid login (cannot be longer than 16 characters)"
    }

    return "0"
}

document.getElementById('login_button').addEventListener('click', function () {
    document.getElementById('login_button').innerHTML = "<span class=\"spinner-border\" role=\"status\" aria-hidden=\"true\"></span>"
    document.getElementById('login_button').disabled = true
    
    var validity = checkLoginValidity(document.getElementById('login_login').value + "|" + document.getElementById('login_password').value)
    if (validity != "0") {
        document.getElementById('login_button').disabled = false
        document.getElementById('login_button').innerHTML = "login"

        validity = validity.split('\\')

        document.getElementById('login_button').classList.remove('my-3')
        document.getElementById('login_response').classList.remove('text-success')
        document.getElementById('login_response').classList.remove('text-danger')
        document.getElementById('login_response').classList.add(validity[0])

        document.getElementById('login_response').innerHTML = validity[1]
        document.getElementById('login_response').style.display = "block"

        return
    }

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
