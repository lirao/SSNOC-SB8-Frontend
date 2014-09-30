var User = require('../models/UserRest');
var Status = require('../models/StatusRest');

module.exports = function(_, io, participants, passport, refreshAllUsers) {
  return {
    getLogin : function(req, res) {
      res.render("join", {message: req.flash('loginMessage')});
    },

    getLogout : function(req, res) {
      req.logout();
      res.redirect('/');
    },

    getSignup : function(req, res) {
        console.log("controller/user.js getSignup");
      res.render('signup', {message: req.flash('signupMessage')});
    },

    getUser : function(req, res) {
      var user_name = req.session.passport.user.user_name;
      User.getUser(user_name, function(err, user) {
        if (user !== null) {
          res.json(200, {name:user.local.name});
        }
      });
    },

    postSignup : function(req, res, next) {
        console.log("controller/user.js postSignup");
      passport.authenticate('local-signup', function(err, user, info) {
        if (err)
          return next(err);
        if (!user)
          return res.redirect('/signup');
        req.logIn(user, function(err) {
          if (err)
            return next(err);
          participants.all.push({'userName' : user.local.name});
          return res.redirect('/welcome');
        });
      })(req, res, next);
    },

    getWelcome : function(req, res) {
      res.render('welcome', {title: "Hello " + req.session.passport.user.user_name + " !!"} );
    },

    postStatus : function(req, res, next) {
        console.log("Posting status for " + req.session.passport.user.user_name);
        user_name = req.session.passport.user.user_name;
        Status.saveNewStatus(user_name, 0, "2014-09-09 09:30", function(err, user) {
            if (user !== null) {
                res.json(200, {name:user.local.name});
            }
            console.log("Error posting status: " + err);
        });
    }


  };
};
