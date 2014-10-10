var online = [];
var offline = [];
var me = {};

function showList(people, myName) {
    var l = people.length;
    for (var i = 0; i < l; i++) {
        if (people[i].name === myName) {
            me = people[i];
            $("#list-me").append("<li id=\"people-me\">"+me.name+" "+me.status,"</li>");
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
        $("#list-online").append("<li class=\"people-online\">"+online[i].name+" "+online[i].status,"</li>");
    }
    offline.sort(compareUser);
    l = offline.length;
    for (var i = 0; i < l; i++) {
        $("#list-offline").append("<li class=\"people-offline\">"+offline[i].name+" "+offline[i].status,"</li>");
    }

}

function updateList(user) {
;
}

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