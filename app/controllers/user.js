var User = require('../models/UserRest');
var Status = require('../models/StatusRest');
var dateFormat = require('dateformat');

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
      console.log("test1");
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

    postStatus : function(req, res) {
        console.log("post");
        var user_name = req.session.passport.user.user_name;
        var user_status = parseInt(req.body.status);
        console.log(user_name);
        console.log(user_status);
        Status.saveNewStatus(user_name, user_status, dateFormat(new Date(), "yyyy-mm-dd HH:MM"), function(err, user) {
            if (user !== null) {
                res.json(200, {name:user.local.name});
            }
            console.log(err);
        });
    }


  };
};
