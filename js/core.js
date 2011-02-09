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

    $("div.sidebar nav ul li a").live("click", function (ev) {
        ev.preventDefault();

        var parent = $(this).parent('li'),
            target = parent.attr('data-target');

        $(this).closest('ul').children('li:not(:has(a))').each(function () {
            var a = document.createElement('A');
            a.href = '#';
            a.innerHTML = $(this).html();
            $(this).html(a);
        });
        parent.html($(this).html());

        $('div.sidebar div.wrapper').addClass('hidden');
        $('div.sidebar div.wrapper.' + target).removeClass('hidden');
    });
});