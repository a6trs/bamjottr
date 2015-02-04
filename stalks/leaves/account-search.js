invite_search_ipt = null;
invite_search_disp = null;
invite_search_prjid = 0;

invite_search_done = function () {
  var resp = JSON.parse(this.responseText);
  if (resp.error !== undefined) {
    invite_search_disp.innerHTML = 'Error: ' + resp.error;
  } else {
    var s;
    // `resp` shoule be an array.
    if (resp.length === 0) {
      s = '<ul class="am-list am-list-static am-list-border">' +
        '<li>No accounts found -~-</li></ul>';
    } else {
      s = '<ul class="am-list am-list-border">';
      resp.forEach(function (a) {
        s += "<li><a href='/invite/" + invite_search_prjid + "/" + a.Account.ID + "'>" + a.Account.Name;
        if (a.Invited) {
          s += " <span class='am-badge am-round am-text-sm badge-invite'>Invited</span>";
        }
        s += '</a></li>';
      });
      s += '</ul>';
    }
    invite_search_disp.innerHTML = s;
  }
};

invite_search_change = function (e) {
  var xhr = new XMLHttpRequest();
  xhr.open('GET', '/account_search/invite/' + invite_search_prjid + '/' + e.target.value);
  xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
  xhr.onload = invite_search_done;
  xhr.send(null);
  // Display the spinner
  invite_search_disp.innerHTML = "<span class='am-icon-spinner am-icon-spin'></span>";
};

invite_search_init = function (prjid) {
  invite_search_ipt = document.getElementById('search_ipt_name');
  invite_search_disp = document.getElementById('search_results');
  invite_search_ipt.onchange = invite_search_change;
  invite_search_prjid = prjid;
};
