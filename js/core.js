$(document).ready(function () {
    $("div.code button").click(function () {
        var options = {},
            code    = $("div.code textarea").val(),
            passed;

        // Get checked options
        $("div.option input:checked").each(function () {
            options[$(this).attr("name")] = true;
        });

        passed = JSHINT(code, options);
    });
});