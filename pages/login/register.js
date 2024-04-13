document.getElementById('register_button').addEventListener('click', function () {
    document.getElementById('register_button').innerHTML = "<span class=\"spinner-border\" role=\"status\" aria-hidden=\"true\"></span>"
    document.getElementById('register_button').disabled = true

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
