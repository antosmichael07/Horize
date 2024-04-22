function checkRegisterValidity(str) {
    var pipeCount = 0
    for (var i = 0; i < str.length; i++) {
        if (str[i] == '|') {
            pipeCount++
        }
    }
    if (pipeCount != 4) {
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

    if (!data[3].toLowerCase().match(/^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|.(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/)) {
        return "text-danger\\Invalid gmail syntax"
    } else if (data[3].split('@')[1] != "gmail.com") {
        return "text-danger\\Invalid gmail (must be a gmail account)"
    }

    return "0"
}

document.getElementById('register_button').addEventListener('click', function () {
    document.getElementById('register_button').innerHTML = "<span class=\"spinner-border\" role=\"status\" aria-hidden=\"true\"></span>"
    document.getElementById('register_button').disabled = true

    var validity = checkRegisterValidity(document.getElementById('register_login').value + "|" + document.getElementById('register_password').value + "|" + document.getElementById('register_nickname').value + "|" + document.getElementById('register_gmail').value + "|" + document.getElementById('register_activation_code').value)
    if (validity != "0") {
        document.getElementById('register_button').disabled = false
        document.getElementById('register_button').innerHTML = "register"

        var validity = validity.split('\\')

        document.getElementById('register_button').classList.remove('my-3')
        document.getElementById('register_response').classList.remove('text-success')
        document.getElementById('register_response').classList.remove('text-danger')
        document.getElementById('register_response').classList.add(validity[0])

        document.getElementById('register_response').innerHTML = validity[1]
        document.getElementById('register_response').style.display = "block"

        return
    }

    fetch("/post_register", {
        method: "POST",
        body: document.getElementById('register_login').value + "|" +
            document.getElementById('register_password').value + "|" +
            document.getElementById('register_nickname').value + "|" +
            document.getElementById('register_gmail').value + "|" +
            document.getElementById('register_activation_code').value
    }).then((response) => {
        response.text().then(function (data) {
            document.getElementById('register_button').disabled = false
            document.getElementById('register_button').innerHTML = "register"

            var res = data.split('\\')

            document.getElementById('register_button').classList.remove('my-3')
            document.getElementById('register_response').classList.remove('text-success')
            document.getElementById('register_response').classList.remove('text-danger')
            document.getElementById('register_response').classList.add(res[0])

            document.getElementById('register_response').innerHTML = res[1]
            document.getElementById('register_response').style.display = "block"
        });
    }).catch((error) => {
        console.log(error)
    });
})
