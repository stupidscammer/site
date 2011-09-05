$.domReady(function () {
  "use strict";

  var templates = {
    error: '<a data-line="%(line)s" href="javascript:void(0)">Line %(line)s</a>: ' +
           '<code>%(code)s</code></p><p>%(msg)s'
  };

  var hasStorage;

  try {
    hasStorage = !!localStorage.getItem && !!JSON;
  } catch (e) {
    hasStorage = false;
  }

  function _(string, context) {
    return string.replace(/%\(\w+\)s/g, function (match) {
      return context[match.slice(2, -2)];
    });
  }

  function listOptions(els, opts) {
    var str = '/*jshint ';

    for (var name in opts) {
      if (opts.hasOwnProperty(name)) {
        str += name + ':' + opts[name] + ', ';
      }
    }

    str = str.slice(0, str.length - 2);
    str += ' */';
    els.append(str);
  }

  function reportFailure(report) {
    var errors = $('div.report ul.jshint-errors');
    var item;

    errors[0].innerHTML = '';
    for (var i = 0, err; err = report.errors[i]; i++) {
      errors.append(_('<li><p>' + templates.error + '</p></li>', {
        line: err.line,
        code: err.evidence,
        msg:  err.reason
      }));

      $('a[data-line="' + err.line + '"]').bind('click', function (ev) {
        var line = $(ev.target).attr('data-line') - 1;
        var str  = Editor.getLine(line);

        Editor.setSelection({ line: line, ch: 0 }, { line: line, ch: str.length });
      });
    }

    listOptions($('div.report > div.error > div.options pre'), report.options);
    $('div.report > div.error').show();
  }

  function reportSuccess(report) {
    listOptions($('div.report > div.success > div.options pre'), report.options);
    $('div.report > div.success').show();
  }

  $('button[data-action=lint]').bind('click', function () {
    var opts   = {};
    var checks = $('ul.inputs-list input[type=checkbox]');

    for (var i = 0, ch; ch = checks[i]; i++) {
      ch = $(ch);

      if (ch.hasClass('neg')) {
        if (!ch.attr('checked')) {
          opts[ch.attr('name')] = true;
        }
      } else {
        if (ch.attr('checked')) {
          opts[ch.attr('name')] = true;
        }
      }
    }

    $('div.report > div.alert-message').hide();
    $('div.report pre').html('');
    JSHINT(Editor.getValue(), opts) ? reportSuccess(JSHINT.data()) : reportFailure(JSHINT.data());
  });

  $('button[data-action=save]').bind('click', function (ev) {
    var button = $(ev.target);
    button.html('Saving report...').attr('disabled', true);

    $('form.save-report textarea[name=code]').val(Editor.getValue());
    $('form.save-report textarea[name=data]').val(JSON.stringify(JSHINT.data()));
    $('form.save-report')[0].submit();
  });

  $('ul.inputs-list input[type=checkbox]').bind('change', function () {
    var checks = $('ul.inputs-list input[type=checkbox]');
    var opts   = {};

    for (var i = 0, ch; ch = checks[i]; i++) {
      ch = $(ch);
      opts[ch.attr('name')] = ch.attr('checked');
    }

    if (hasStorage) {
      localStorage.setItem('opts', JSON.stringify(opts));
    }
  });

  if (!hasStorage) {
    return;
  }

  var checks = $('ul.inputs-list input[type=checkbox]');
  var opts   = JSON.parse(localStorage.getItem('opts') || '{}');

  for (var i = 0, ch; ch = checks[i]; i++) {
    ch = $(ch);
    if (opts[ch.attr('name')] === true) {
      ch.attr('checked', true);
    } else if (opts[ch.attr('name')] === false) {
      ch.removeAttr('checked');
    }
  }
});