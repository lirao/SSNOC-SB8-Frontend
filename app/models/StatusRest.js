/**
 * Created by ltanady on 29/9/14.
 */
var request = require('request');
var rest_api = require('../../config/rest_api');

function Status(user_name, status_code) {
    this.local = {
        name : user_name,
        status : status_code
    }

}

Status.saveNewStatus = function(user_name, status_code, create_at, callback) {
    var options = {
        url : rest_api.post_new_status + user_name,
        body : {userName: user_name, statusCode: status_code, createdAt: create_at},
        json : true
    }
    console.log("models/StatusRest.js saveNewStatus");
    request.post(options, function(err, res, body) {
        if (err) {
            callback(err, null);
            return;
        }
        if (res.statusCode !== 200 && res.statusCode !== 201) {
            callback(res.body, null);
            return;
        }
        console.log("Response from create status: " + res.body);
        var new_status = new Status(user_name, status_code);
        callback(null, new_status);
        return;
    });
}

module.exports = Status;