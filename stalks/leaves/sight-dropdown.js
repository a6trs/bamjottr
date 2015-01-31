sight_items = [];
sight_dropdown = document.getElementById('sight-btn');
sight_lastsel = 0;

disp_sel_sight_item = function (idx) {
  'use strict';
  sight_dropdown.innerHTML = sight_items[idx].innerHTML + "&nbsp;<span class='am-icon-caret-down'></span>";
  sight_items[sight_lastsel].classList.remove('am-active');
  sight_items[idx].classList.add('am-active');
  sight_lastsel = idx;
};

init_sight_dropdown = function (prjid) {
  'use strict';

  var sight_selctd = function (idx) {
    if (idx === sight_lastsel) return;
    // stackoverflow.com/q/247483
    // stackoverflow.com/q/9713058
    var xhr = new XMLHttpRequest();
    xhr.open('POST', '/sight');
    xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
    xhr.setRequestHeader('Connection', 'close');
    // Target type: 0 = project, 1 = post
    // See [grass/sights.go] for details
    var params = 'tgttype=0&tgtid=' + prjid + '&level=' + idx;
    xhr.send(params);
    disp_sel_sight_item(idx);
  };

  var handler = function (idx) {
    return function () { sight_selctd(idx); };
  };

  var i;
  for (i = 0; i < 3; i++) {
    sight_items[i] = document.getElementById('sight-sel-' + i.toString());
    sight_items[i].onclick = handler(i);
  }

};
