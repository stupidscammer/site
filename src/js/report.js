$.domReady(function () {
  "use strict";

  var Editor = CodeMirror.fromTextArea(document.getElementById('editor'), {
    theme: 'default',
    lineNumbers: true,
    readOnly: true
  });

  var el = Editor.getScrollerElement();
  el.style.height = "600px";
  Editor.refresh();

  var templates = {
    error: '<a class="goto" data-line="%(line)s" href="javascript:void(0)">Line %(line)s</a>: ' +
           '<code>%(code)s</code></p><p>%(msg)s'
  };

  var optarr = $('.options > pre').text().replace('/*jshint ', '').replace('* /', '').split(', ');
  var opts = {};
  var name, val;

  function _(string, context) {
    return string.replace(/%\(\w+\)s/g, function (match) {
      return context[match.slice(2, -2)];
    });
  }

  for (var i = 0, opt; opt = optarr[i]; i++) {
    name = opt.split(':')[0];
    val  = opt.split(':')[1];

    if (val == 'true' || val == 'false') {
      opts[name] = (val == 'true');
    } else {
      opts[name] = val;
    }
  }

  JSHINT( $('#editor').val(), opts);
  var errors = $('div.report ul.jshint-errors');
  var report = JSHINT.data();
  var item;

  errors[0].innerHTML = '';
  for (var i = 0, err; err = report.errors[i]; i++) {
    errors.append(_('<li><p>' + templates.error + '</p></li>', {
      line: err.line,
      code: err.evidence,
      msg:  err.reason
    }));
  }

  $('a.goto').bind('click', function (ev) {
      var line = $(ev.target).attr('data-line') - 1;
      var str  = Editor.getLine(line);

      Editor.setSelection({ line: line, ch: 0 }, { line: line, ch: str.length });
  });
});