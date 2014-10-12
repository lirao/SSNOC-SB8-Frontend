var online = [];
var offline = [];
var me = {};
var sock = new WebSocket('ws://localhost:3000/status');
sock.onopen = function(){
    sock.onmessage = function(res){
        var body = JSON.parse(res.data);
        switch (body.type) {
        case "people":
            showList(body.users);
            break;
        case "status":
            updateStatus(body.user);
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
    sock.send(JSON.stringify({"type": "people"}));
}

function showList(people) {
    var myName = $('#me-name').text();
    var l = people.length;
    for (var i = 0; i < l; i++) {
        user = {"name": people[i].userName, "status": people[i].lastStatusCode.toString(), "online": "1"};
        if (people[i].userName === myName) {
            me = user;
            $('#status').val(me.status);
            changeColor($('#me-status'), me.status);
        } else {
                online.push(user);
//            if (people[i].online === 1) {
//                online.push(user);
//            } else {
//                offline.push(user);
//            }
        }
    }
    online.sort(compareUser);
    l = online.length;
    for (var i = 0; i < l; i++) {
        $("#list-online").append("<li class=\"list-group-item\">"+online[i].name+"<b>&#9679;</b></li>");
        changeColor($("#list-online").children().eq(i).find("b"), online[i].status);
    }
    offline.sort(compareUser);
    l = offline.length;
    for (var i = 0; i < l; i++) {
        $("#list-offline").append("<li class=\"list-group-item\">"+offline[i].name+"<b>&#9679;</b></li>");
        changeColor($("#list-offline").children().eq(i).find("b"), online[i].status);
    }
}

function updateStatus(user) {
    if (user.name === me.name) {
        return;
    }
    var onlineIndex = findUserByName(online, user.name);
    online[onlineIndex] = user;
    changeColor($("#list-online").children().eq(onlineIndex).find("b"), user.status);

//    var offlineIndex = findUserByName(offline, user.name);
//    if (offlineIndex !== -1) {
//        if (user.online === "1") {
//            offline.splice(offlineIndex, 1);
//            online.push(user);
//            online.sort(compareUser);
//            var onlineIndex = findUserByName(online, user.name);
//            if (onlineIndex + 1 === online.length) {
//                $("#list-online").append($("#list-offline").children().eq(offlineIndex));
//                changeColor($("#list-online").last().find("b"), user.status);
//            } else {
//                $("#list-online").children().eq(onlineIndex).before($("#list-offline").children().eq(offlineIndex));
//                changeColor($("#list-online").children().eq(onlineIndex).find("b"), user.status);
//            }
//        } else {
//            offline[offlineIndex] = user;
//            changeColor($("#list-offline").children().eq(offlineIndex).find("b"), user.status);
//        }
//    } else {
//        var onlineIndex = findUserByName(online, user.name);
//        if (onlineIndex !== -1) {
//        } else {
//
//        }
//    }


//    if (user.online === 1) { // New online user
//        var offlineIndex = findUserByName(offline, user.name);
//        if (offlineIndex !== -1) { // Old offline user
//            offline.splice(offlineIndex, 1);
//            online.push(user);
//            online.sort(compareUser);
//            var onlineIndex = findUserByName(online, user.name);
//            if (onlineIndex + 1 === online.length) {
//                $("#list-online").append($("#list-offline").children().eq(offlineIndex));
//            } else {
//                $("#list-online").children().eq(onlineIndex).before($("#list-offline").children().eq(offlineIndex));
//            }
//        } else { // New user
//            online.push(user);
//            online.sort(compareUser);
//            var onlineIndex = findUserByName(online, user.name);
//            if (onlineIndex + 1 === online.length) {
//                $("#list-online").append("<li class=\"list-group-item\">"+user.name+"<b>&#9679;</b></li>");
//            } else {
//                $("#list-online").children().eq(onlineIndex).before("<li class=\"list-group-item\">"+user.name+"<b>&#9679;</b></li>");
//            }
//        }
//    } else { // New offline user
//        var onlineIndex = findUserByName(online, user.name);
//        if (onlineIndex !== -1) { // Old online user
//            online.splice(onlineIndex, 1);
//            offline.push(user);
//            offline.sort(compareUser);
//            var offlineIndex = findUserByName(offline, user.name);
//            if (offlineIndex + 1 === offline.length) {
//                $("#list-offline").append($("#list-online").children().eq(onlineIndex));
//            } else {
//                $("#list-offline").children().eq(offlineIndex).before($("#list-online").children().eq(onlineIndex));
//            }
//        } else { // New user
//            offline.push(user);
//            offline.sort(compareUser);
//            var offlineIndex = findUserByName(offline, user.name);
//            if (offlineIndex + 1 === offline.length) {
//                $("#list-offline").append("<li class=\"list-group-item\">"+user.name+"<b>&#9679;</b></li>");
//            } else {
//                $("#list-offline").children().eq(offlineIndex).before("<li class=\"list-group-item\">"+user.name+"<b>&#9679;</b></li>");
//            }
//        }
//    }
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