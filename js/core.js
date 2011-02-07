$(document).ready(function () {
    $("div.code button").click(function () {
        alert(JSHINT($("div.code textarea").val()));
    });
});