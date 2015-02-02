sight_items = [];
sight_dropdown = document.getElementById('sight-btn');
sight_lastsel = -1;

// Target type: 0 = project, 1 = post
// See [grass/sights.go] for details
sight_target = {
  PROJECT: 0,
  POST: 1
};

var remove_last = function (s, pattern) {
  return s.substr(0, s.lastIndexOf(pattern));
};

var inc_text = function (el, inc) {
  el.innerHTML = (parseInt(el.innerHTML) + inc).toString();
};

disp_sel_sight_item = function (idx) {
  'use strict';
  // Remove contents after the last '&nbsp;'
  sight_dropdown.innerHTML = remove_last(sight_items[idx].innerHTML, '&nbsp;') + "&nbsp;<span class='am-icon-caret-down'></span>";
  if (sight_lastsel !== -1) {
    // Decrease the count of the last selected item by 1
    sight_items[sight_lastsel].parentElement.classList.remove('am-active');
    inc_text(sight_items[sight_lastsel].getElementsByClassName('am-badge')[0], -1);
    inc_text(sight_items[idx].getElementsByClassName('am-badge')[0], 1);
  }
  sight_items[idx].parentElement.classList.add('am-active');
  sight_lastsel = idx;
};

init_sight_dropdown = function (targetType, targetId) {
  'use strict';

  var sight_selctd = function (idx) {
    if (idx === sight_lastsel) return;
    // stackoverflow.com/q/247483
    // stackoverflow.com/q/9713058
    var xhr = new XMLHttpRequest();
    xhr.open('POST', '/sight');
    xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
    xhr.setRequestHeader('Connection', 'close');
    var params = 'tgttype=' + targetType + '&tgtid=' + targetId + '&level=' + idx;
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
