var online = [];
var offline = [];
var me = {};

function showList(people, myName) {
    var l = people.length;
    for (var i = 0; i < l; i++) {
        if (people[i].name === myName) {
            me = people[i];
            $("#list-me").append("<li>"+me.name+" "+me.status,"</li>");
        } else {
            if (people[i].online === 1) {
                online.push(people[i]);
            } else {
                offline.push(people[i]);
            }
        }
    }
    online.sort(compareUser);
    l = online.length;
    for (var i = 0; i < l; i++) {
        $("#list-online").append("<li>"+online[i].name+" "+online[i].status,"</li>");
    }
    offline.sort(compareUser);
    l = offline.length;
    for (var i = 0; i < l; i++) {
        $("#list-offline").append("<li>"+offline[i].name+" "+offline[i].status,"</li>");
    }
}

function updateList(user) {
    if (user.online === 1) { // New online user
        var offlineIndex = findUserByName(offline, user.name);
        if (offlineIndex !== -1) { // Old offline user
            offline.splice(offlineIndex, 1);
            online.push(user);
            online.sort(compareUser);
            var onlineIndex = findUserByName(online, user.name);
            if (onlineIndex + 1 === online.length) {
                $("#list-online").append($("#list-offline").children().eq(offlineIndex));
            } else {
                $("#list-online").children().eq(onlineIndex).before($("#list-offline").children().eq(offlineIndex));
            }
        } else { // New user
            online.push(user);
            online.sort(compareUser);
            var onlineIndex = findUserByName(online, user.name);
            if (onlineIndex + 1 === online.length) {
                $("#list-online").append("<li>"+user.name+" "+user.status,"</li>");
            } else {
                $("#list-online").children().eq(onlineIndex).before("<li>"+user.name+" "+user.status,"</li>");
            }
        }
    } else { // New offline user
        var onlineIndex = findUserByName(online, user.name);
        if (onlineIndex !== -1) { // Old online user
            online.splice(onlineIndex, 1);
            offline.push(user);
            offline.sort(compareUser);
            var offlineIndex = findUserByName(offline, user.name);
            if (offlineIndex + 1 === offline.length) {
                $("#list-offline").append($("#list-online").children().eq(onlineIndex));
            } else {
                $("#list-offline").children().eq(offlineIndex).before($("#list-online").children().eq(onlineIndex));
            }
        } else { // New user
            offline.push(user);
            offline.sort(compareUser);
            var offlineIndex = findUserByName(offline, user.name);
            if (offlineIndex + 1 === offline.length) {
                $("#list-offline").append("<li>"+user.name+" "+user.status,"</li>");
            } else {
                $("#list-offline").children().eq(offlineIndex).before("<li>"+user.name+" "+user.status,"</li>");
            }
        }
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

var people = [{'name':"me", 'status':0, 'online':1},
              {'name':"online-user1",'status':1,'online':1},
              {'name':"online-user2",'status':2,'online':1},
              {'name':"online-user0",'status':0,'online':1},
              {'name':"offline-user2",'status':2,'online':0},
              {'name':"offline-user0",'status':0,'online':0},
              {'name':"offline-user1",'status':1,'online':0}];

$( document ).ready(showList(people, "me"));

$("#test1").click(function() {updateList({'name':"offline-user1",'status':1,'online':1})});
$("#test2").click(function() {updateList({'name':"xnew-user",'status':55,'online':1})});