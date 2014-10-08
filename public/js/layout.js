$(document).ready(function() {
    var pathname = window.location.pathname === "/" ? "welcome" : window.location.pathname.substring(1);
    var navid = "#nav-" + pathname;
    if (navid === "#nav-welcome" || navid === "#nav-people") {
        $(navid).addClass("active");
    }
    if (navid === "#nav-login" || navid === "#nav-signup") {
        $(".nav-off").show();
        $(".nav-on").hide();
    } else {
        $(".nav-on").show();
        $(".nav-off").hide();
    }
});