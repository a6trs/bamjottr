invite_search_ipt = null;
invite_search_disp = null;

invite_search_done = function () {
  invite_search_disp.innerHTML = this.responseText;
};

invite_search_change = function (e) {
  var xhr = new XMLHttpRequest();
  // Testing. Add account search handler and remove this test.
  xhr.open('GET', '/login');//'/account_search?type=invite&q=' + e.target.value);
  xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
  xhr.onload = invite_search_done;
  xhr.send(null);
};

invite_search_init = function (prjid) {
  invite_search_ipt = document.getElementById('search_ipt_name');
  invite_search_disp = document.getElementById('search_results');
  invite_search_ipt.onchange = invite_search_change;
};
