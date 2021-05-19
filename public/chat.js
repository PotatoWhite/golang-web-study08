$(function () {
    if (!window.EventSource) {
        alert("No Event Source");
        return;
    }

    let $chatLog = $('#chat-log');
    let $chatMsg = $('#chat-msg');

    let isBlank = function (string) {
        return string == null || string.trim() === "";
    }

    let userName;
    while (isBlank(userName)) {
        userName = prompt("What's your name?");
        if (!isBlank(userName))
            $('#user-name').html('<b>' + userName + '</b>');
    }

    $('#input-form').on('submit', function (e) {
        $.post('/messages', {
            msg: $chatMsg.val(),
            name: userName
        });

        $chatMsg.val("");
        $chatMsg.focus();
        return false;
    });

    let addMessage = function (data) {
        let text = "";
        if (!isBlank(data.name)) {
            text = '<strong>' + data.name + ':</strong>';
        }

        text += data.msg;
        $chatLog.prepend('<div><span>' + text + '</span></div>');
    }

    let es = new EventSource('/stream');
    es.onopen = function (e) {
        $.post('/users', {
            name: userName
        })
    };

    es.onmessage = function (e) {
        message = JSON.parse(e.data);
        console.log(message)
        addMessage(message);
    };

    window.onbeforeunload = function () {
        es.close();
        alert("bye")
        $.ajax({
            url: "/users?name=" + userName,
            type: "DELETE"
        });
    };
})