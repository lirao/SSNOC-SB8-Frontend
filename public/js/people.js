var online = [];
var offline = [];
var me = {};
var sock = new WebSocket('ws://localhost:3000/socket');
sock.onopen = function(){
    sock.onmessage = function(res){
        var body = JSON.parse(res.data);
        switch (body.type) {
        case "people":
            updatePeople(body.users);
            break;
        case "status":
            updateStatus(body.user);
            break;
        case "login":
            updateLogin(body.user);
            break;
        case "logout":
            updateLogout(body.user);
            break;
        case "signup":
            updateSignup(body.user);
            break;
        case "wall":
            updateWall(body.messages);
            break;
        case "chat":
            postMessage(body.message);
            break;
        }
        var newMessage = $('<li>').text(res.data);
            $('#messages').append(newMessage);
            $('#message').val('');
    };
    $('#status').on('change', function() {
        me.status = $('#status').val();
        changeColor($('#me-status'), me.status);
        sock.send(JSON.stringify({"user": me, "type": "status"}));
        return false;
    });
    $('#post-wall').on('click', function() {
        if ($('input[name="wall-message"]').val().length > 0) {
            sock.send(JSON.stringify({"type": "chat", "message": {"sender": me.name, "receiver": "public", "content": $('input[name="wall-message"]').val()}}));
            $('input[name="wall-message"]').val('');
            $('input[name="wall-message"]').focus();
        }
    });
    sock.send(JSON.stringify({"type": "people"}));
    sock.send(JSON.stringify({"type": "wall"}));
}

function updatePeople(people) {
    var myName = $('#me-name').text();
    for (var temp in people) {
        user = {"name": people[temp].Name, "status": people[temp].Status, "online": people[temp].Online};
        if (people[temp].Name === myName) {
            me = user;
            $('#status').val(me.status);
            changeColor($('#me-status'), me.status);
        } else {
            if (people[temp].Online === "1") {
                online.push(user);
            } else {
                offline.push(user);
            }
        }
    }
    online.sort(compareUser);
    l = online.length;
    for (var i = 0; i < l; i++) {
        $("#list-online").append("<li class=\"list-group-item\">"+online[i].name+"<b>&#9679;</b></li>");
        changeColor($("#list-online").children().eq(i+1).find("b"), online[i].status);
    }
    offline.sort(compareUser);
    l = offline.length;
    for (var i = 0; i < l; i++) {
        $("#list-offline").append("<li class=\"list-group-item\">"+offline[i].name+"<b>&#9679;</b></li>");
        changeColor($("#list-offline").children().eq(i+1).find("b"), offline[i].status);
    }
}

function updateStatus(user) {
    if (user.name === me.name) {
        return;
    }
    var onlineIndex = findUserByName(online, user.name);
    online[onlineIndex] = user;
    changeColor($("#list-online").children().eq(onlineIndex+1).find("b"), user.status);
}

function updateLogin(user) {
    if (user.Name === me.name) {
        return;
    }
    var offlineIndex = findUserByName(offline, user.Name);
    offline.splice(offlineIndex, 1);
    temp = {"name": user.Name, "status": user.Status, "online": user.Online};
    online.push(temp);
    online.sort(compareUser);
    var onlineIndex = findUserByName(online, temp.name);
    if (onlineIndex + 1 === online.length) {
        $("#list-online").append($("#list-offline").children().eq(offlineIndex + 1));
    } else {
        $("#list-online").children().eq(onlineIndex + 1).before($("#list-offline").children().eq(offlineIndex + 1));
    }
}

function updateLogout(user) {
    if (user.Name === me.name) {
        return;
    }
    var onlineIndex = findUserByName(online, user.Name);
    online.splice(onlineIndex, 1);
    temp = {"name": user.Name, "status": user.Status, "online": user.Online};
    offline.push(temp);
    offline.sort(compareUser);
    var offlineIndex = findUserByName(offline, temp.name);
    if (offlineIndex + 1 === offline.length) {
        $("#list-offline").append($("#list-online").children().eq(onlineIndex + 1));
    } else {
        $("#list-offline").children().eq(offlineIndex + 1).before($("#list-online").children().eq(onlineIndex + 1));
    }
}

function updateSignup(user) {
    if (user.Name === me.name) {
        return;
    }
    temp = {"name": user.Name, "status": user.Status, "online": user.Online};
    online.push(temp);
    online.sort(compareUser);
    var onlineIndex = findUserByName(online, temp.name);
    if (onlineIndex + 1 === online.length) {
        $("#list-online").append("<li class=\"list-group-item\">"+temp.name+"<b>&#9679;</b></li>");
    } else {
        $("#list-online").children().eq(onlineIndex + 1).before("<li class=\"list-group-item\">"+temp.name+"<b>&#9679;</b></li>");
    }
    changeColor($("#list-online").children().eq(onlineIndex+1).find("b"), temp.status);
}

function updateWall(messages) {
    var l = messages.length;
    for (var i = 0; i < l; i++) {
        postMessage(messages[i]);
    }
}

function postMessage(message) {
    if (message.receiver === "public") {
        $('#wall').children().eq(0).after("<li class=\"list-group-item\">"+message.sender+"("+message.time+"): "+message.content+"</li>");
    }
}

function findUserByName(list, name){
    var l = list.length;
    for (var i = 0; i < l; i++) {
        if (list[i].name === name) {
            return i;
        }
    }
    return -1;
};

function compareUser(user0, user1) {
    return user0.name.localeCompare(user1.name);
}

function changeColor(elem, status) {
    switch (status) {
    case "0":
        elem.css({"color":"green"});
        break;
    case "1":
        elem.css({"color":"yellow"});
        break;
    case "2":
        elem.css({"color":"red"});
        break;
    }
    return;
}