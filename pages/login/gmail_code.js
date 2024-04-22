function checkGmailValidity(str) {
    if (!str.toLowerCase().match(/^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|.(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/)) {
        return "text-danger\\Invalid gmail syntax"
    } else if (str.split('@')[1] != "gmail.com") {
        return "text-danger\\Invalid gmail (must be a gmail account)"
    }

    return "0"
}

document.getElementById('send_code_button').addEventListener('click', function () {
    document.getElementById('send_code_button').innerHTML = "<span class=\"spinner-border spinner-border-sm mt-1\" role=\"status\" aria-hidden=\"true\"></span>"
    document.getElementById('send_code_button').disabled = true

    var validity = checkGmailValidity(document.getElementById('register_gmail').value)
    if (validity != "0") {
        document.getElementById('send_code_button').disabled = false
        document.getElementById('send_code_button').innerHTML = "send code"

        var validity = validity.split('\\')

        document.getElementById('send_code_button').classList.remove('mb-4')
        document.getElementById('register_gmail').classList.remove('mb-4')
        document.getElementById('send_code_button').classList.add('mb-0')
        document.getElementById('register_gmail').classList.add('mb-0')

        document.getElementById('gmail_code_response').classList.remove('text-success')
        document.getElementById('gmail_code_response').classList.remove('text-danger')
        document.getElementById('gmail_code_response').classList.add(validity[0])

        document.getElementById('gmail_code_response').innerHTML = validity[1]
        document.getElementById('gmail_code_response').style.display = "block"

        return
    }

    fetch("/post_gmail_code", {
        method: "POST",
        body: document.getElementById('register_gmail').value
    }).then((response) => {
        response.text().then(function (data) {
            document.getElementById('send_code_button').disabled = false
            document.getElementById('send_code_button').innerHTML = "send code"

            var res = data.split('\\')

            document.getElementById('send_code_button').classList.remove('mb-4')
            document.getElementById('register_gmail').classList.remove('mb-4')
            document.getElementById('send_code_button').classList.add('mb-0')
            document.getElementById('register_gmail').classList.add('mb-0')

            document.getElementById('gmail_code_response').classList.remove('text-success')
            document.getElementById('gmail_code_response').classList.remove('text-danger')
            document.getElementById('gmail_code_response').classList.add(res[0])

            document.getElementById('gmail_code_response').innerHTML = res[1]
            document.getElementById('gmail_code_response').style.display = "block"
        });
    }).catch((error) => {
        console.log(error)
    });
})
