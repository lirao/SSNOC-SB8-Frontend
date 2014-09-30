var User=require('../models/UserRest');

module.exports = function(_, io, participants, passport) {
  return {
    getPeople: function(req, res) {
        var userName = req.session.passport.user.user_name;
        var userStatus;
        participants.all = [];
        User.getAllUsers(function(err, users) {
            if (!err) {
                users.forEach(function (user) {
                    if (user.local.name === userName) {
                        userStatus = user.local.status;
                    }
                    participants.all.push({
                        userName: user.local.name,
                        userStatus: user.local.status
                    });
                });
                switch (userStatus) {
                    case 0:
                        userStatus = "/img/green-dot.png";
                        break;
                    case 1:
                        userStatus = "/img/yellow-dot.png";
                        break;
                    case 2:
                        userStatus = "/img/red-dot.png";
                        break;
                    case 3:
                        userStatus = "/img/grey-dot.png";
                        break;

                }
                res.render("people", {userId: req.session.userId, title:"People", user_name:userName, user_status:userStatus});
            }
        });
    }
  };
};