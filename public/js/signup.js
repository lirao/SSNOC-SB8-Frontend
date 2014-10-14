$('#signup').on('submit', function() {
    if ($('input[name="userName"]').val().length < 3) {
        $('#signup').find("font:first").text("User name requires at least 3 characters.");
        $('input[name="password"]').val('');
        $('input[name="passwordRe"]').val('');
        $('input[name="userName"]').focus();
        return false;
    }
    if ($('input[name="password"]').val().length < 4) {
        $('#signup').find("font:first").text("Password requires at least 4 characters.");
        $('input[name="password"]').val('');
        $('input[name="passwordRe"]').val('');
        $('input[name="password"]').focus();
        return false;
    }
    if ($('input[name="password"]').val() !== $('input[name="passwordRe"]').val()) {
        $('#signup').find("font:first").text("Two passwords should be the same.");
        $('input[name="password"]').val('');
        $('input[name="passwordRe"]').val('');
        $('input[name="password"]').focus();
        return false;
    }
    return true;
});