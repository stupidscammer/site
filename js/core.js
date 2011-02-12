$(document).ready(function () {
    function onPass() {
        var output   = $("div.output.wrapper"),
            template = $("#passOutput");

        output.removeClass("fail");
        output.addClass("pass");
        output.html(template.tmpl(JSHINT.data()));
    }

    function onFail() {
        var output   = $("div.output.wrapper"),
            template = $("#failOutput");

        output.removeClass("pass");
        output.addClass("fail");
        output.html(template.tmpl(JSHINT.data()));
    }

    function onEmpty() {
        var output   = $("div.output.wrapper"),
            template = $("#emptyOutput");

        output.html(template.tmpl());
    }

    $("div.code button").click(function () {
        var options = {},
            code    = $("div.code textarea").val(),
            output  = $("div.sidebar li[data-target=output] a");

        // Get checked options
        $("div.option input:checked").each(function () {
            options[$(this).attr("name")] = true;
        });

        if (jQuery.trim(code) === "")
            onEmpty();
        else if (JSHINT(code, options))
            onPass();
        else
            onFail();

        if (output)
            output.trigger("click");
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