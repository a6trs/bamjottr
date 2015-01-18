checkform_signup = function (e) {
  var p1 = document.getElementById('ipt-pwd').value;
  var p2 = document.getElementById('ipt-pwd-rpt').value;
  var hinter = document.getElementById('errmsg');
  var errmsgs = [];
  if (p1.length < 6) {
    errmsgs.push("Password too short");
  } else if (p1 !== p2) {
    errmsgs.push("Passwords doesn't match :(");
  }
  if (errmsgs.length > 0) {
    if (hinter.innerHTML === '') hinter.classList.add('bg-danger');
    hinter.innerHTML = errmsgs.join('<br>');
    return false;
  }
  return true;
};
