invite_search_ipt = null;
invite_search_disp = null;
invite_search_prjid = 0;

invite_search_done = function () {
  var resp = JSON.parse(this.responseText);
  if (resp.error !== undefined) {
    invite_search_disp.innerHTML = 'Error: ' + resp.error;
  } else {
    s = '';
    // `resp` shoule be an array.
    resp.forEach(function (a) {
      s += "<a href='/invite/" + invite_search_prjid + "/" + a.Account.ID + "'>" + a.Account.Name + "</a>";
      if (a.Invited) {
        s += " <span class='am-badge am-round badge-invite'>Invited</span>";
      }
      s += '<br>';
    });
    invite_search_disp.innerHTML = s;
  }
};

invite_search_change = function (e) {
  var xhr = new XMLHttpRequest();
  xhr.open('GET', '/account_search/invite/' + invite_search_prjid + '/' + e.target.value);
  xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
  xhr.onload = invite_search_done;
  xhr.send(null);
};

invite_search_init = function (prjid) {
  invite_search_ipt = document.getElementById('search_ipt_name');
  invite_search_disp = document.getElementById('search_results');
  invite_search_ipt.onchange = invite_search_change;
  invite_search_prjid = prjid;
};
